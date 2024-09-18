package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

const (
	caPubKey = `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMH00QoOliUpItlPrJ3QOd6DRH+wgd0vDb1k5FZxO9Iy CA`
	hostKey  = `-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----
`
)

// ─ ssh -F sshconfig -J jumphost targethost
func main() {
	server := &ssh.Server{
		Addr: ":2222",
		ChannelHandlers: map[string]ssh.ChannelHandler{
			"direct-tcpip": directTCPIPHandler,
		},
		PublicKeyHandler: certificateAuthenticator,
	}

	s, err := gossh.ParsePrivateKey([]byte(hostKey))
	if err != nil {
		log.Fatal(err)
	}
	server.AddHostKey(s)

	log.Println("Starting SSH jump server on :2222")
	log.Fatal(server.ListenAndServe())
}

func directTCPIPHandler(srv *ssh.Server, conn *gossh.ServerConn, newChan gossh.NewChannel, ctx ssh.Context) {
	d := struct {
		DestAddr string
		DestPort uint32
		SrcAddr  string
		SrcPort  uint32
	}{}

	if err := gossh.Unmarshal(newChan.ExtraData(), &d); err != nil {
		newChan.Reject(gossh.ConnectionFailed, "Failed to parse channel data")
		return
	}

	dest := fmt.Sprintf("%s:%d", d.DestAddr, d.DestPort)

	var dialer net.Dialer
	dconn, err := dialer.Dial("tcp", dest)
	if err != nil {
		newChan.Reject(gossh.ConnectionFailed, fmt.Sprintf("Failed to connect to %s: %v", dest, err))
		return
	}

	ch, reqs, err := newChan.Accept()
	if err != nil {
		dconn.Close()
		return
	}
	go gossh.DiscardRequests(reqs)

	go func() {
		defer ch.Close()
		defer dconn.Close()
		io.Copy(ch, dconn)
	}()
	go func() {
		defer ch.Close()
		defer dconn.Close()
		io.Copy(dconn, ch)
	}()

	log.Printf("Proxying connection from %s:%d to %s:%d", d.SrcAddr, d.SrcPort, d.DestAddr, d.DestPort)
}

func certificateAuthenticator(ctx ssh.Context, key ssh.PublicKey) bool {
	trustedCAKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(caPubKey))
	if err != nil {
		log.Fatalf("Failed to parse CA public key: %v", err)
	}

	checker := gossh.CertChecker{
		IsUserAuthority: func(auth gossh.PublicKey) bool {
			// Verify if the given public key is your trusted CA
			return bytes.Equal(auth.Marshal(), trustedCAKey.Marshal())
		},
	}

	cert, ok := key.(*gossh.Certificate)
	if !ok {
		log.Printf("Is certificate: %t", ok)
		return false
	}

	if cert.CertType != gossh.UserCert {
		log.Print("CertType is not of UserCert")
		return false
	}

	if !checker.IsUserAuthority(cert.SignatureKey) {
		log.Print("Certificate not signed by CA")
		return false
	}

	if err := checker.CheckCert(ctx.User(), cert); err != nil {
		log.Printf("CheckCert failed: %v", err)
		return false
	}

	return true
}
