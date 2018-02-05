package paypal

const (
	K_PAYPAL_EVENT_RESOURCE_TYPE_INVOICES = "invoices"
	K_PAYPAL_EVENT_RESOURCE_TYPE_SALE     = "sale"
)

// https://developer.paypal.com/docs/integration/direct/webhooks/event-names/
const (
	K_PAYPAL_EVENT_TYPE_PAYMENT_SALE_COMPLETED = "PAYMENT.SALE.COMPLETED"
	K_PAYPAL_EVENT_TYPE_PAYMENT_SALE_REFUNDED  = "PAYMENT.SALE.REFUNDED"
)

type EventType struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

type Webhook struct {
	Id         string       `json:"id,omitempty"`
	URL        string       `json:"url"`
	EventTypes []*EventType `json:"event_types,omitempty"`
	Links      []*Link      `json:"links,omitempty"`
}

type WebhookList struct {
	Webhooks []*Webhook `json:"webhooks,omitempty"`
}

type Event struct {
	Id           string      `json:"id"`
	CreateTime   string      `json:"create_time,omitempty"`
	ResourceType string      `json:"resource_type,omitempty"`
	EventVersion string      `json:"event_version,omitempty"`
	EventType    string      `json:"event_type,omitempty"`
	Summary      string      `json:"summary,omitempty"`
	Resource     interface{} `json:"resource,omitempty"`
	Status       string      `json:"status,omitempty"`
	Links        []*Link     `json:"links,omitempty"`
}

type verifyWebhookSignatureParam struct {
	AuthAlgo         string `json:"auth_algo"`
	CertURL          string `json:"cert_url"`
	TransmissionId   string `json:"transmission_id"`
	TransmissionSig  string `json:"transmission_sig"`
	TransmissionTime string `json:"transmission_time"`
	WebhookId        string `json:"webhook_id"`
	WebhookEvent     *Event `json:"webhook_event"`
}

type verifyWebhookSignatureResponse struct {
	VerificationStatus string `json:"verification_status"`
}
