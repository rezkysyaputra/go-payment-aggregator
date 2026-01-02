package pkg

func MapMidtransStatus(midtransStatus, fraudStatus string) string {
	switch midtransStatus {
	case "capture":
		if fraudStatus == "challenge" {
			return "PENDING"
		}
		return "PAID"
	case "settlement":
		return "PAID"
	case "pending":
		return "PENDING"
	case "deny", "failure", "cancel":
		return "FAILED"
	case "expire":
		return "FAILED"
	default:
		return "PENDING"
	}
}
