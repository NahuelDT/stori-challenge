package email

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
	"time"

	"github.com/tartoide/stori/stori-challenge/internal/domain"
)

const emailTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Your Storicard Monthly Snapshot</title>
    <style>
        body {
            font-family: 'Inter', Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #333333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #F0F0F5;
        }
        .container {
            background-color: #FFFFFF;
            padding: 35px;
            border-radius: 12px;
            box-shadow: 0 4px 15px rgba(0,0,0,0.08);
        }
        .header {
            text-align: center;
            border-bottom: 1px solid #E0E0E0;
            padding-bottom: 25px;
            margin-bottom: 30px;
        }
        .logo img {
            width: 100%;
            display: block;
            height: auto;
            max-height: 220px;
            margin-bottom: 15px;
            object-fit: cover;
        }
        .title {
            font-size: 26px;
            font-weight: 600;
            color: #0E5858;
            margin-bottom: 10px;
        }
        .greeting {
            font-size: 16px;
            color: #555555;
            margin-bottom: 30px;
            text-align: center;
        }
        .balance-section {
            text-align: center;
            margin-bottom: 35px;
        }
        .balance-label {
            font-size: 16px;
            color: #555555;
            margin-bottom: 8px;
        }
        .balance-amount {
            font-size: 36px;
            font-weight: 700;
            color: #0E5858;
            padding: 15px;
            background-color: #F8F9FA;
            border-radius: 8px;
            display: inline-block;
        }
        .section {
            margin: 30px 0;
        }
        .section-title {
            font-size: 20px;
            font-weight: 600;
            color: #0E5858;
            margin-bottom: 20px;
            border-bottom: 1px solid #DCDCDC;
            padding-bottom: 12px;
        }
        .transaction-item, .average-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 16px 8px;
            border-bottom: 1px solid #E9E9E9;
            font-size: 15px;
        }
        .transaction-item:last-child, .average-item:last-child {
            border-bottom: none;
        }
        .item-label {
            color: #4A4A4A;
            margin-right: 10px;
        }
        .item-value {
            font-weight: 700;
            color: #0E5858;
        }
        .credit {
            color: #1A7A7A;
            font-weight: 700;
        }
        .debit {
            color: #D32F2F;
            font-weight: 700;
        }
        .footer {
            text-align: center;
            margin-top: 40px;
            padding-top: 25px;
            border-top: 1px solid #E0E0E0;
            color: #777777;
            font-size: 13px;
        }
        .footer p {
            margin: 5px 0;
        }
        .footer a {
            color: #0E5858;
            text-decoration: none;
        }
        .footer a:hover {
            text-decoration: underline;
        }
        .no-transactions {
            text-align: center;
            color: #666;
            font-style: italic;
            padding: 25px;
            background-color: #F8F9FA;
            border-radius: 8px;
            margin: 30px 0;
        }
        @media screen and (max-width: 480px) {
            .container {
                padding: 20px;
            }
            .title {
                font-size: 22px;
            }
            .balance-amount {
                font-size: 28px;
            }
            .section-title {
                font-size: 18px;
            }
            .transaction-item, .average-item {
                font-size: 14px;
                padding: 12px 4px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">
                <img src="https://stori-challenge-nahue.s3.sa-east-1.amazonaws.com/stori-bg-1.avif" alt="Storicard Logo">
            </div>
            <div class="title">Your Monthly Stori Snapshot</div>
        </div>

        <div class="greeting">
            Hi! Here's a summary of your account activity this month.
        </div>
        
        <div class="balance-section">
            <div class="balance-label">Your Current Balance</div>
            <div class="balance-amount">${{.TotalBalance.StringFixed 2}}</div>
        </div>
        
        {{if .HasTransactions}}
        <div class="section">
            <div class="section-title">📈 Monthly Transactions Overview</div>
            {{range .SortedMonthlyTransactions}}
            <div class="transaction-item">
                <span class="item-label">Transactions in {{.Month}}:</span>
                <span class="item-value"><strong>{{.Count}}</strong></span>
            </div>
            {{end}}
        </div>
        
        <div class="section">
            <div class="section-title">📊 Average Transaction Values</div>
            {{if not .AverageCredit.IsZero}}
            <div class="average-item">
                <span class="item-label">Average credit:</span>
                <span class="credit">+${{.AverageCredit.StringFixed 2}}</span>
            </div>
            {{end}}
            {{if not .AverageDebit.IsZero}}
            <div class="average-item">
                <span class="item-label">Average debit:</span>
                <span class="debit">-${{.AverageDebit.StringFixed 2}}</span>
            </div>
            {{end}}
            {{if .AverageCredit.IsZero}} {{if .AverageDebit.IsZero}}
            <div class="average-item">
                <span class="item-label">No credit or debit transactions this month.</span>
            </div>
            {{end}}{{end}}
        </div>
        {{else}}
        <div class="no-transactions">
            🌱 Looks like there were no transactions for your account in the processed file this month.
        </div>
        {{end}}
        
        <div class="footer">
            <p>Thank you for being a Storicard member!</p>
            <p>Have questions? Visit our <a href="https://www.storicard.com/faq">FAQ</a> or <a href="https://www.storicard.com/contact">contact support</a>.</p>
            <p>This is an automated message. Please do not reply directly to this email.</p>
            <p>&copy; {{ .CurrentYear }} Storicard. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

type templateData struct {
	*domain.Summary
	SortedMonthlyTransactions []monthTransaction
	CurrentYear               int
}

type monthTransaction struct {
	Month string
	Count int
}

// RenderEmailTemplate renders the email template with summary data
func RenderEmailTemplate(summary *domain.Summary) (string, error) {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing email template: %w", err)
	}

	var sortedMonthly []monthTransaction
	for month, count := range summary.MonthlyTransactions {
		sortedMonthly = append(sortedMonthly, monthTransaction{
			Month: month,
			Count: count,
		})
	}

	sort.Slice(sortedMonthly, func(i, j int) bool {
		return sortedMonthly[i].Month < sortedMonthly[j].Month
	})

	data := templateData{
		Summary:                   summary,
		SortedMonthlyTransactions: sortedMonthly,
		CurrentYear:               time.Now().Year(),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing email template: %w", err)
	}

	return buf.String(), nil
}
