package network

import (
	"context"
	"crypto/sha512"
	tls "crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cossacklabs/acra/network/testutils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/credentials"
	"net"
	"runtime/debug"
	"strings"
	"sync"
	"testing"
	"time"
)

func getConnectionPair(address string, listener net.Listener, t testing.TB) (net.Conn, net.Conn) {
	serverConnCh := make(chan net.Conn)
	clientConnCh := make(chan net.Conn)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}
		serverConnCh <- conn
	}()
	go func() {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			t.Fatal(err)
		}
		clientConnCh <- conn
	}()
	// wait when client connect to server
	var clientConn, serverConn net.Conn
	for i := 0; i < 2; i++ {
		select {
		case clientConn = <-clientConnCh:
			continue
		case serverConn = <-serverConnCh:
			continue
		case <-time.NewTimer(time.Second / 2).C:
			t.Fatal("Timeout on connection client with server")
		}
	}
	return clientConn, serverConn
}

func getTLSConfigs(t testing.TB) (*tls.Config, *tls.Config) {
	serverTLSConfig, err := NewTLSConfig("localhost", "", "", "", tls.RequireAndVerifyClientCert, NewCertVerifierAll())
	if err != nil {
		t.Fatal(err)
	}
	clientTLSConfig, err := NewTLSConfig("localhost", "", "", "", tls.RequireAndVerifyClientCert, NewCertVerifierAll())
	if err != nil {
		t.Fatal(err)
	}
	clientConfig, serverConfig, err := testutils.GetTestTLSConfigs(func() *tls.Config { return clientTLSConfig }, func() *tls.Config { return serverTLSConfig })
	if err != nil {
		t.Fatal(err)
	}
	return clientConfig, serverConfig
}

func TestTLSWrapperWithCertificateAuthentication(t *testing.T) {
	expectedClientCommonName := []byte("CN=client1,OU=IT,O=Global Security,L=London,C=GB")
	value := sha512.Sum512(expectedClientCommonName)
	expectedClientID := []byte(hex.EncodeToString(value[:]))
	clientConfig, serverConfig := getTLSConfigs(t)
	converter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	extractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper, err := NewTLSAuthenticationConnectionWrapper(true, nil, serverConfig, extractor)
	if err != nil {
		t.Fatal(err)
	}
	clientWrapper, err := NewTLSAuthenticationConnectionWrapper(false, clientConfig, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	testWrapper(clientWrapper, serverWrapper, expectedClientID, wrapperCommunicationIterations, t)
}

func TestTLSWrapperWithStaticClientID(t *testing.T) {
	clientID := []byte(`some client id`)
	clientConfig, serverConfig := getTLSConfigs(t)
	serverWrapper, err := NewTLSConnectionWrapper(clientID, serverConfig)
	if err != nil {
		t.Fatal(err)
	}
	clientWrapper, err := NewTLSConnectionWrapper(nil, clientConfig)
	if err != nil {
		t.Fatal(err)
	}
	testWrapper(clientWrapper, serverWrapper, clientID, wrapperCommunicationIterations, t)
}

func BenchmarkTLSWrapper(t *testing.B) {
	//  openssl x509 -in client1.crt -subject -noout -nameopt RFC2253 | sed 's/subject=//'
	expectedClientCommonName := []byte("CN=client1,OU=IT,O=Global Security,L=London,ST=London,C=GB")
	expectedClientID := []byte(hex.EncodeToString(expectedClientCommonName))
	clientConfig, serverConfig := getTLSConfigs(t)
	converter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	extractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper, err := NewTLSAuthenticationConnectionWrapper(true, nil, serverConfig, extractor)
	if err != nil {
		t.Fatal(err)
	}
	clientWrapper, err := NewTLSAuthenticationConnectionWrapper(true, clientConfig, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	testWrapper(clientWrapper, serverWrapper, expectedClientID, t.N, t)
}

// isTLS13 return true if connection has version > tls12 constant value in ConnectionState after successful handshake
func isTLS13(conn net.Conn) bool {
	// check with GREATER comparison because golang versions < 1.2 have not constant VersionTLS13
	return UnwrapSafeCloseConnection(conn).(*tls.Conn).ConnectionState().Version > tls.VersionTLS12
}

func TestTLSConfigWeakCipherSuitDeny(t *testing.T) {
	clientConfig, serverConfig := getTLSConfigs(t)
	converter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	extractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper, err := NewTLSAuthenticationConnectionWrapper(true, nil, serverConfig, extractor)
	if err != nil {
		t.Fatal(err)
	}
	clientWrapper, err := NewTLSAuthenticationConnectionWrapper(false, clientConfig, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	const host = "localhost"
	listener, err := net.Listen("tcp", host+":0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	clientConn, serverConn := getConnectionPair(fmt.Sprintf("%s:%d", host, port), listener, t)

	wrapErrorCh := make(chan bool)
	// check not allowed cipher suit
	clientWrapper.clientConfig.CipherSuites = []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256}

	go func() {
		conn, _, err := serverWrapper.WrapServer(context.TODO(), serverConn)
		if err != nil {
			if err.Error() != "tls: no cipher suite supported by both client and server" {
				t.Fatal("Expected error with unsupported ciphersuits")
			}
			wrapErrorCh <- true
			return
		}
		// tls1.3 in golang doesn't support ciphersuites configuration, so just return ok
		if isTLS13(conn) {
			wrapErrorCh <- true
			return
		}
		t.Fatal("expected error")
	}()
	go func() {
		conn, err := clientWrapper.WrapClient(context.TODO(), clientConn)
		if err != nil {
			if err.Error() != "remote error: tls: handshake failure" {
				t.Fatal("Expected with handshake failure")
			}
			wrapErrorCh <- true
			return
		}
		// tls1.3 in golang doesn't support ciphersuites configuration, so just return ok
		if isTLS13(conn) {
			wrapErrorCh <- true
			return
		}
		t.Fatal("expected error")
	}()
	for i := 0; i < 2; i++ {
		select {
		case <-wrapErrorCh:
			continue
		case <-time.NewTimer(time.Second / 2).C:
			t.Fatal("Timeout on wrap with incorrect cipher suits")
		}
	}
	if err = clientConn.Close(); err != nil {
		t.Fatal(err)
	}
	if err = serverConn.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestTLSConfigWeakVersion(t *testing.T) {
	clientConfig, serverConfig := getTLSConfigs(t)
	clientID := []byte(`some client`)
	serverWrapper, err := NewTLSConnectionWrapper(clientID, serverConfig)
	if err != nil {
		t.Fatal(err)
	}
	clientWrapper, err := NewTLSConnectionWrapper(nil, clientConfig)
	if err != nil {
		t.Fatal(err)
	}
	clientWrapper.clientConfig.MinVersion = tls.VersionSSL30
	clientWrapper.clientConfig.MaxVersion = tls.VersionTLS11

	matchedServerSideError := func(err error) bool {
		expectedMessages := []string{
			// go < 1.12
			"tls: client offered an unsupported, maximum protocol version of",
			// go >= 1.12
			"tls: client offered only unsupported versions"}
		found := false
		for _, msg := range expectedMessages {
			if strings.HasPrefix(err.Error(), msg) {
				found = true
			}
		}
		return found
	}

	matchedClientSideError := func(err error) bool {
		return err.Error() == "remote error: tls: protocol version not supported"
	}

	matchedServerSide := false
	matchedClientSide := false
	mutex := sync.Mutex{}
	// we expects 2 errors, one from client side and from server side related with protocol version is unsupported
	onError := func(err error, t testing.TB) {
		mutex.Lock()
		defer mutex.Unlock()
		if matchedServerSide && matchedClientSide {
			return
		}
		if matchedClientSideError(err) {
			matchedClientSide = true
		}
		if matchedServerSideError(err) {
			matchedServerSide = true
		}
		if !(matchedServerSide || matchedClientSide) {
			t.Fatalf("Unexpected error %s\n", err)
		}
	}
	testWrapperWithError(clientWrapper, serverWrapper, clientID, 1, onError, t)
}
func TestTLSCertificateAuthenticationByCommonName(t *testing.T) {
	clientConfig, serverConfig := getTLSConfigs(t)
	//  openssl x509 -in client1.crt -subject -noout -nameopt RFC2253 | sed 's/subject=//'
	expectedClientCommonName := []byte("CN=client1,OU=IT,O=Global Security,L=London,C=GB")
	value := sha512.Sum512(expectedClientCommonName)
	expectedClientID := []byte(hex.EncodeToString(value[:]))
	serverHost := "localhost"
	clientConfig.ServerName = serverHost

	converter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	extractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper, err := NewTLSAuthenticationConnectionWrapper(true, nil, serverConfig, extractor)
	if err != nil {
		t.Fatal(err)
	}
	clientWrapper, err := NewTLSAuthenticationConnectionWrapper(false, clientConfig, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	testWrapper(clientWrapper, serverWrapper, expectedClientID, wrapperCommunicationIterations, t)
}

func TestTLSCertificateAuthenticationBySerialNumber(t *testing.T) {
	clientConfig, serverConfig := getTLSConfigs(t)
	converter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	extractor, err := NewTLSClientIDExtractor(SerialNumberExtractor{}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper, err := NewTLSAuthenticationConnectionWrapper(true, nil, serverConfig, extractor)
	if err != nil {
		t.Fatal(err)
	}
	clientWrapper, err := NewTLSAuthenticationConnectionWrapper(false, clientConfig, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	clientCertificate, err := x509.ParseCertificate(clientConfig.Certificates[0].Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	value := sha512.Sum512(clientCertificate.SerialNumber.Bytes())
	expectedClientID := []byte(hex.EncodeToString(value[:]))
	testWrapper(clientWrapper, serverWrapper, expectedClientID, wrapperCommunicationIterations, t)
}

func TestEmptyCertificateChain(t *testing.T) {
	clientConfig, serverConfig := getTLSConfigs(t)
	converter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	extractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper, err := NewTLSAuthenticationConnectionWrapper(true, nil, serverConfig, extractor)
	if err != nil {
		t.Fatal(err)
	}
	// remove client's CA to not pass verification to check that empty VerifiedChain not pass
	serverWrapper.serverConfig.ClientCAs = x509.NewCertPool()
	serverWrapper.serverConfig.ClientAuth = tls.RequireAnyClientCert
	clientWrapper, err := NewTLSAuthenticationConnectionWrapper(false, clientConfig, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	clientCertificate, err := x509.ParseCertificate(clientConfig.Certificates[0].Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	expectedClientID, err := getClientIDFromCertificate(clientCertificate, extractor)
	if err != nil {
		t.Fatal(err)
	}
	// expect that first error will be ErrNoPeerCertificate
	tested := false
	mutex := sync.Mutex{}
	onError := func(err error, t testing.TB) {
		mutex.Lock()
		defer mutex.Unlock()
		if tested {
			return
		}
		if err == ErrNoPeerCertificate {
			tested = true
			return
		}
		t.Fatalf("Expected error ErrNoPeerCertificate, took %s\n", err)
	}
	testWrapperWithError(clientWrapper, serverWrapper, expectedClientID, 1, onError, t)
}

func TestClientsCertificateDenyOnValidation(t *testing.T) {
	clientConfig, serverConfig := getTLSConfigs(t)
	ca := generateTLSCA(t)
	clientCertificateTemplate := generateCertificateTemplate(t)
	clientCertificateTemplate.IsCA = true
	clientCertificateTemplate.Subject.CommonName = "client"
	clientCertificate := createLeafKey(ca, clientCertificateTemplate, t)
	clientConfig.Certificates = []tls.Certificate{clientCertificate}
	caCrt, err := x509.ParseCertificate(ca.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	serverConfig.ClientCAs.AddCert(caCrt)

	clientWrapper, err := NewTLSAuthenticationConnectionWrapper(false, clientConfig, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	converter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	extractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper, err := NewTLSAuthenticationConnectionWrapper(true, nil, serverConfig, extractor)
	if err != nil {
		t.Fatal(err)
	}
	clientCrt, err := x509.ParseCertificate(clientCertificate.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	expectedClientID, err := extractor.ExtractClientID(clientCrt)
	if err != nil {
		t.Fatal(err)
	}
	// expect that first error will be ErrCACertificateUsed
	tested := false
	mutex := sync.Mutex{}
	onError := func(err error, t testing.TB) {
		mutex.Lock()
		defer mutex.Unlock()
		if tested {
			return
		}
		if err == ErrCACertificateUsed {
			tested = true
			return
		}
		// need to print stack because hard to clarify where exactly was error in communication of client with server side
		debug.PrintStack()
		t.Fatalf("Expected error ErrCACertificateUsed, took %s\n", err)
	}
	testWrapperWithError(clientWrapper, serverWrapper, expectedClientID, 1, onError, t)
}

type testExtractor struct{ err error }

func (e testExtractor) GetCertificateIdentifier(certificate *x509.Certificate) ([]byte, error) {
	return nil, e.err
}

func TestClientsCertificateDenyOnClientIDExtraction(t *testing.T) {
	clientConfig, serverConfig := getTLSConfigs(t)
	clientWrapper, err := NewTLSAuthenticationConnectionWrapper(false, clientConfig, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	converter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	extractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper, err := NewTLSAuthenticationConnectionWrapper(true, nil, serverConfig, extractor)
	if err != nil {
		t.Fatal(err)
	}

	clientCrt, err := x509.ParseCertificate(clientConfig.Certificates[0].Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	expectedClientID, err := getClientIDFromCertificate(clientCrt, extractor)
	if err != nil {
		t.Fatal(err)
	}

	// override extractor which will always return err
	expectedErr := errors.New("test error")
	testExtractor, err := NewTLSClientIDExtractor(testExtractor{err: expectedErr}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper.clientIDExtractor = testExtractor

	tested := false
	mutex := sync.Mutex{}
	onError := func(err error, t testing.TB) {
		mutex.Lock()
		defer mutex.Unlock()
		if tested {
			return
		}
		if err == expectedErr {
			tested = true
			return
		}
		t.Fatalf("Expected error expectedErr, took %s\n", err)
	}
	testWrapperWithError(clientWrapper, serverWrapper, expectedClientID, 1, onError, t)
}

type testConvertor struct{ err error }

func (t testConvertor) Convert(identifier []byte) ([]byte, error) {
	return nil, t.err
}

func TestClientsCertificateDenyOnClientIDConvertation(t *testing.T) {
	clientConfig, serverConfig := getTLSConfigs(t)
	clientWrapper, err := NewTLSAuthenticationConnectionWrapper(false, clientConfig, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	converter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	extractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, converter)
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper, err := NewTLSAuthenticationConnectionWrapper(true, nil, serverConfig, extractor)
	if err != nil {
		t.Fatal(err)
	}

	clientCrt, err := x509.ParseCertificate(clientConfig.Certificates[0].Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	expectedClientID, err := getClientIDFromCertificate(clientCrt, extractor)
	if err != nil {
		t.Fatal(err)
	}

	// override convertor which will always return err
	expectedErr := errors.New("test error")
	testExtractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, testConvertor{err: expectedErr})
	if err != nil {
		t.Fatal(err)
	}
	serverWrapper.clientIDExtractor = testExtractor

	tested := false
	mutex := sync.Mutex{}
	onError := func(err error, t testing.TB) {
		mutex.Lock()
		defer mutex.Unlock()
		if tested {
			return
		}
		if err == expectedErr {
			tested = true
			return
		}
		t.Fatalf("Expected error expectedErr, took %s\n", err)
	}
	testWrapperWithError(clientWrapper, serverWrapper, expectedClientID, 1, onError, t)
}

func generateTLSCA(t testing.TB) tls.Certificate {
	ca, err := testutils.GenerateTLSCA()
	if err != nil {
		t.Fatal(err)
	}
	return ca
}

func generateTLSCAFromTemplate(caTemplate *x509.Certificate, t testing.TB) tls.Certificate {
	ca, err := testutils.GenerateTLSCAFromTemplate(caTemplate)
	if err != nil {
		t.Fatal(err)
	}
	return ca
}
func createLeafKey(caCert tls.Certificate, templateCertificate *x509.Certificate, t testing.TB) tls.Certificate {
	leafKey, err := testutils.CreateLeafKey(caCert, templateCertificate)
	if err != nil {
		t.Fatal(err)
	}
	return leafKey
}
func generateCertificateTemplate(t testing.TB) *x509.Certificate {
	template, err := testutils.GenerateCertificateTemplate()
	if err != nil {
		t.Fatal(err)
	}
	return template

}

func TestTLSGRPCClientIDExtractorSuccess(t *testing.T) {
	clientCert := generateCertificateTemplate(t)
	testClientID := "client1"
	clientCert.Subject.CommonName = testClientID
	idConverter, err := NewDefaultHexIdentifierConverter()
	if err != nil {
		t.Fatal(err)
	}
	tlsClientIDExtractor, err := NewTLSClientIDExtractor(DistinguishedNameExtractor{}, idConverter)
	if err != nil {
		t.Fatal(err)
	}
	expectedClientID, err := tlsClientIDExtractor.ExtractClientID(clientCert)
	if err != nil {
		t.Fatal(err)
	}
	authInfo := &wrappedTLSAuthInfo{
		TLSInfo: credentials.TLSInfo{State: tls.ConnectionState{VerifiedChains: [][]*x509.Certificate{{clientCert}}}},
		conn:    newClientIDConnection(&testConnection{}, []byte(expectedClientID)),
	}
	resultClientID, err := GetClientIDFromAuthInfo(authInfo, tlsClientIDExtractor)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, resultClientID, expectedClientID)
}

func TestTLSGRPCClientIDExtractorIncorrectAuthInfo(t *testing.T) {
	testTLSGRPCClientIDExtractorIncorrectAuthInfo(t)
}

func TestClientAuth0WithClientIDExtraction(t *testing.T) {
	_, err := NewTLSAuthenticationConnectionWrapper(true, nil, &tls.Config{}, nil)
	if err != ErrInvalidTLSConfiguration {
		t.Fatalf("Expect ErrInvalidTLSConfiguration, took %s\n", err.Error())
	}
}
