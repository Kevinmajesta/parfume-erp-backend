package email

import (
	"fmt"

	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Config *entity.Config
}

func NewEmailSender(config *entity.Config) *EmailSender {
	return &EmailSender{Config: config}
}

func (e *EmailSender) SendEmail(to []string, subject, body string) error {
	from := "kepvivv@gmail.com"
	password := e.Config.SMTP.Password
	smtpHost := e.Config.SMTP.Host

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(smtpHost, 587, from, password)
	err := dialer.DialAndSend(mailer)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (e *EmailSender) SendWelcomeEmail(to, name, message string) error {
	subject := "Welcome Email | Depublic"
	body := fmt.Sprintf("Dear %s,\nThis is a welcome email message from depublic\n\nDepublic Team", name)
	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailSender) SendResetPasswordEmail(to, name, resetCode string) error {
	subject := "Reset Password | Depublic"
	body := fmt.Sprintf("Dear %s,\nPlease use the following code to reset your password: %s\n\nDepublic Team", name, resetCode)
	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailSender) SendVerificationEmail(to, name, code string) error {
	subject := "Verify Your Email | Depublic"
	body := fmt.Sprintf("Dear %s,\nPlease use the following code to verify your email: %s\n\nDepublic Team", name, code)
	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailSender) SendTransactionInfo(to, Transactions_id, Cart_id, User_id,
	Fullname_user, Trx_date, Payment, Payment_url, Amount string) error {
	subject := "Transaction Info | Depublic"
	body := fmt.Sprintf("Dear %s,\nThis is your transaction info from Depublic:\n\nTransaction ID: %s\n\nCart ID: %s\n\nUser ID: %s\n\nFullname: %s\n\nTransaction Date: %s\n\nPayment Type: %s\n\nURL Payment: %s\n\nTotal Amount: %s\n\n\nDepublic Team ",
		Fullname_user, Transactions_id, Cart_id, User_id, Fullname_user, Trx_date, Payment, Payment_url, Amount)
	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailSender) SendRfqEmail(to, rfqId, vendorId, orderDate, status string, products []entity.RfqsProduct) error {
	subject := fmt.Sprintf("RFQ Confirmation | RFQ ID: %s", rfqId)

	// HTML body untuk email
	body := fmt.Sprintf(`
    <html>
    <head>
        <style>
            body {
                font-family: Arial, sans-serif;
                color: #333;
                background-color: #f4f4f4;
                margin: 0;
                padding: 0;
            }
            .container {
                width: 80%;
                margin: auto;
                background-color: #fff;
                padding: 20px;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            }
            h2 {
                color: #2c3e50;
            }
            table {
                width: 100%;
                border-collapse: collapse;
                margin: 20px 0;
            }
            table, th, td {
                border: 1px solid #ddd;
                padding: 8px;
            }
            th {
                background-color: #f2f2f2;
                text-align: left;
            }
            tr:nth-child(even) {
                background-color: #f9f9f9;
            }
            footer {
                text-align: center;
                margin-top: 30px;
                font-size: 12px;
                color: #777;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h2>RFQ Details | RFQ ID: %s</h2>
            <p><strong>Dear Vendor,</strong></p>
            <p>Here are the details of your RFQ:</p>

            <table>
                <tr>
                    <th>RFQ ID</th>
                    <td>%s</td>
                </tr>
                <tr>
                    <th>Vendor ID</th>
                    <td>%s</td>
                </tr>
                <tr>
                    <th>Order Date</th>
                    <td>%s</td>
                </tr>
                <tr>
                    <th>Status</th>
                    <td>%s</td>
                </tr>
            </table>

            <h3>Products:</h3>
            <table border="1" cellpadding="5" cellspacing="0" style="border-collapse: collapse; width: 100%; margin-bottom: 20px; border: 1px solid #ddd;">
        <tr>
            <th style="background-color: #f2f2f2; text-align: left; padding: 8px;">Product Name</th>
            <th style="background-color: #f2f2f2; text-align: left; padding: 8px;">Quantity</th>
            <th style="background-color: #f2f2f2; text-align: left; padding: 8px;">Unit Price</th>
            <th style="background-color: #f2f2f2; text-align: left; padding: 8px;">Subtotal</th>
        </tr>`,
		rfqId, rfqId, rfqId, rfqId, vendorId, orderDate, status) // Memperbaiki urutan parameter

	for _, product := range products {
		body += fmt.Sprintf(`
				<tr>
					<td style="padding: 8px;">%s</td>
					<td style="padding: 8px; text-align: center;">%s</td> <!-- Format sebagai string -->
					<td style="padding: 8px; text-align: right;">%s</td> <!-- Format sebagai string -->
					<td style="padding: 8px; text-align: right;">%s</td> <!-- Format sebagai string -->
				</tr>`,
			product.ProductName,
			fmt.Sprintf("%s", product.Quantity),  // Konversi Quantity ke string
			fmt.Sprintf("%s", product.UnitPrice), // Konversi UnitPrice ke string
			fmt.Sprintf("%s", product.Subtotal))  // Konversi Subtotal ke string
	}

	// Menutup tabel dan body email
	body += fmt.Sprintf(`
            </table>

            <footer>
                <p>Thank you for doing business with us!</p>
                <p>Best regards,<br>Depublic Team</p>
            </footer>
        </div>
    </body>
    </html>
`)

	// Kirim email
	return e.SendEmail([]string{to}, subject, body)

}

func (e *EmailSender) SendQuoEmail(to, quotationId, costumerId, orderDate, status string, products []entity.QuotationsProduct) error {
	subject := fmt.Sprintf("Quotation Confirmation | Quotation ID: %s", quotationId)

	// HTML body untuk email
	body := fmt.Sprintf(`
    <html>
    <head>
        <style>
            body {
                font-family: Arial, sans-serif;
                color: #333;
                background-color: #f4f4f4;
                margin: 0;
                padding: 0;
            }
            .container {
                width: 80%;
                margin: auto;
                background-color: #fff;
                padding: 20px;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            }
            h2 {
                color: #2c3e50;
            }
            table {
                width: 100%;
                border-collapse: collapse;
                margin: 20px 0;
            }
            table, th, td {
                border: 1px solid #ddd;
                padding: 8px;
            }
            th {
                background-color: #f2f2f2;
                text-align: left;
            }
            tr:nth-child(even) {
                background-color: #f9f9f9;
            }
            footer {
                text-align: center;
                margin-top: 30px;
                font-size: 12px;
                color: #777;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h2>Quotation Details | Quotation ID: %s</h2>
            <p><strong>Dear Costumer,</strong></p>
            <p>Here are the details of your Quotation:</p>

            <table>
                <tr>
                    <th>Quotation ID</th>
                    <td>%s</td>
                </tr>
                <tr>
                    <th>Costumer ID</th>
                    <td>%s</td>
                </tr>
                <tr>
                    <th>Order Date</th>
                    <td>%s</td>
                </tr>
                <tr>
                    <th>Status</th>
                    <td>%s</td>
                </tr>
            </table>

            <h3>Products:</h3>
            <table border="1" cellpadding="5" cellspacing="0" style="border-collapse: collapse; width: 100%; margin-bottom: 20px; border: 1px solid #ddd;">
        <tr>
            <th style="background-color: #f2f2f2; text-align: left; padding: 8px;">Product Name</th>
            <th style="background-color: #f2f2f2; text-align: left; padding: 8px;">Quantity</th>
            <th style="background-color: #f2f2f2; text-align: left; padding: 8px;">Unit Price</th>
            <th style="background-color: #f2f2f2; text-align: left; padding: 8px;">Subtotal</th>
        </tr>`,
		quotationId, quotationId, quotationId, quotationId, costumerId, orderDate, status) // Memperbaiki urutan parameter

	for _, product := range products {
		body += fmt.Sprintf(`
				<tr>
					<td style="padding: 8px;">%s</td>
					<td style="padding: 8px; text-align: center;">%s</td> <!-- Format sebagai string -->
					<td style="padding: 8px; text-align: right;">%s</td> <!-- Format sebagai string -->
					<td style="padding: 8px; text-align: right;">%s</td> <!-- Format sebagai string -->
				</tr>`,
			product.ProductName,
			fmt.Sprintf("%s", product.Quantity),  // Konversi Quantity ke string
			fmt.Sprintf("%s", product.UnitPrice), // Konversi UnitPrice ke string
			fmt.Sprintf("%s", product.Subtotal))  // Konversi Subtotal ke string
	}

	// Menutup tabel dan body email
	body += fmt.Sprintf(`
            </table>

            <footer>
                <p>Thank you for doing business with us!</p>
                <p>Best regards,<br>Depublic Team</p>
            </footer>
        </div>
    </body>
    </html>
`)

	// Kirim email
	return e.SendEmail([]string{to}, subject, body)

}
