package ldap

//TODO remove asn1 dependancies here
import (
	"encoding/asn1"
	"fmt"

	"github.com/vjeantet/asn1-ber"
)

type ProtocolOp interface {
	String() string
}

type LDAPRequest interface {
	GetMessageId() int
	String() string
	GetProtocolOp() ProtocolOp
}

type Message struct {
	wroteMessage int
	MessageId    int
	ProtocolOp   ProtocolOp
	Controls     []interface{}
	client       *client
}

func (m Message) GetMessageId() int {
	return m.MessageId
}

func (m Message) String() string {
	return fmt.Sprintf("MessageId=%d, %s", m.MessageId, m.ProtocolOp.String())
}
func (m Message) GetProtocolOp() ProtocolOp {
	return m.ProtocolOp
}

func (m *Message) SetClient(c *client) {
	m.client = c
}

func (m *Message) GetClient() *client {
	return m.client
}

// BIND REQUEST MESSAGE
type BindRequest struct {
	Message
	ProtocolOp struct {
		Version  int
		Login    []byte
		Password []byte
	}
}

func (b *BindRequest) SetLogin(login []byte) {
	b.ProtocolOp.Login = login
}

func (b *BindRequest) GetLogin() []byte {
	return b.ProtocolOp.Login
}

func (b *BindRequest) SetVersion(version int) {
	b.ProtocolOp.Version = version
}

func (b *BindRequest) SetPassword(password []byte) {
	b.ProtocolOp.Password = password
}

func (b *BindRequest) GetPassword() []byte {
	return b.ProtocolOp.Password
}

func (b BindRequest) String() string {
	var s string = ""

	s = fmt.Sprintf("Login:%s, Password:%s",
		b.GetLogin(),
		b.GetPassword())

	return s
}

// UNBIND REQUEST MESSAGE
type UnbindRequest struct {
	Message
	ProtocolOp struct {
	}
}

// SEARCH REQUEST MESSAGE
type SearchRequest struct {
	Message
	ProtocolOp struct {
		BaseDN       []byte
		Scope        int
		DerefAliases int
		SizeLimit    int
		TimeLimit    int
		TypesOnly    bool
		FilterRaw    asn1.RawValue
		Attributes   [][]byte
		Filter       string
	}
}

func (s *SearchRequest) GetTypesOnly() bool {
	return s.ProtocolOp.TypesOnly
}
func (s *SearchRequest) GetFilterRaw() asn1.RawValue {
	return s.ProtocolOp.FilterRaw
}
func (s *SearchRequest) GetAttributes() [][]byte {
	return s.ProtocolOp.Attributes
}
func (s *SearchRequest) GetFilter() string {
	return s.ProtocolOp.Filter
}
func (s *SearchRequest) GetBaseDN() []byte {
	return s.ProtocolOp.BaseDN
}
func (s *SearchRequest) GetScope() int {
	return s.ProtocolOp.Scope
}
func (s *SearchRequest) GetDerefAliases() int {
	return s.ProtocolOp.DerefAliases
}
func (s *SearchRequest) GetSizeLimit() int {
	return s.ProtocolOp.SizeLimit
}
func (s *SearchRequest) GetTimeLimit() int {
	return s.ProtocolOp.TimeLimit
}

func (r SearchRequest) String() string {
	var s string = ""

	s = fmt.Sprintf("BaseDn:%s\nScope:%d\nDerefAliases:%d\nSizeLimit:%d\nTimeLimit:%d\nTypesOnly:%t\nFilter:%s\n",
		r.ProtocolOp.BaseDN,
		r.ProtocolOp.Scope,
		r.ProtocolOp.DerefAliases,
		r.ProtocolOp.SizeLimit,
		r.ProtocolOp.TimeLimit,
		r.ProtocolOp.TypesOnly,
		r.ProtocolOp.Filter)

	for i := range r.ProtocolOp.Attributes {
		s = fmt.Sprintf("%sAttribute:%s\n", s, r.ProtocolOp.Attributes[i])
	}

	return s
}

// REPONSES
type LDAPResponse interface {
	encodeToAsn1() []byte
}
type LDAPResult struct {
	ResultCode        int
	MatchedDN         string
	DiagnosticMessage string
	referral          interface{}
}

// BindResponse
type BindResponse struct {
	LDAPResult
	Request         *BindRequest
	serverSaslCreds string
}

func (b *BindResponse) Send() {
	b.Request.GetClient().write(b.encodeToAsn1())
	b.Request.wroteMessage += 1
}

func (b BindResponse) encodeToAsn1() []byte {
	packet := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSequence, nil, "LDAP Response")
	packet.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimative, ber.TagInteger, uint64(b.Request.GetMessageId()), "MessageID"))

	bindResponse := ber.Encode(ber.ClassApplication, ber.TypeConstructed, ApplicationBindResponse, nil, "Bind Response")
	bindResponse.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimative, ber.TagEnumerated, uint64(b.ResultCode), "ResultCode"))
	bindResponse.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimative, ber.TagOctetString, b.MatchedDN, "MatchedDN"))
	bindResponse.AppendChild(ber.NewString(ber.ClassUniversal, ber.TypePrimative, ber.TagOctetString, b.DiagnosticMessage, "DiagnosticMessage"))

	packet.AppendChild(bindResponse)
	return packet.Bytes()
}

func (r BindResponse) String() string {
	return ""
}

type SearchResponse struct {
	LDAPResult

	Entries   []*Entry
	Referrals []string
	//Controls []Control
	chan_out chan Message
}

func (sr *SearchResponse) AddEntry(entry *Entry) {
	sr.Entries = append(sr.Entries, entry)
}

func (r SearchResponse) String() string {
	return ""
}

type Entry struct {
	DN         string
	Attributes []*EntryAttribute
}

func (e *Entry) SetDn(dn string) {
	e.DN = dn
}

func (e *Entry) AddAttribute(name string, values ...string) {
	var ea = &EntryAttribute{Name: name, Values: values}
	e.Attributes = append(e.Attributes, ea)
}

type EntryAttribute struct {
	Name   string
	Values []string
}

type UnbindResponse struct {
	LDAPResult
}

func (msg UnbindResponse) String() string {
	return ""
}