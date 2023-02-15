# iupp itau promotion scan
### [EN]
This is a simple project written in Go that runs a Github Action every day to check for new promotions to exchange iupp points for brazilian airline miles (TudoAzul / Smiles / Latam Pass) and send an email to the email list using [SendGrid](https://sendgrid.com/) if a new promotion shows up.
If you want to receive the alert email as well you can fork this project and setup your own github secrets.

Required secrets:
```
SENDGRID_API_KEY
SENDGRID_SENDER_EMAIL
SENDGRID_TO_EMAIL
SENDGRID_TO_NAME
```

### [PT/BR]
Este é um projeto escrito em Go que executa uma Github Action todos os dias para verificar novas promoções para trocar pontos iupp por milhas aéreas (TudoAzul / Smiles / Latam Pass) e envia um email usando [SendGrid](https://sendgrid.com/) caso uma nova promoção apareça.
Se você deseja receber um e-mail de promoção, você pode fazer um fork desse projeto e setar seus próprios github secrets.

Secrets necessários:
```
SENDGRID_API_KEY
SENDGRID_SENDER_EMAIL
SENDGRID_TO_EMAIL
SENDGRID_TO_NAME
```

#Contrib
Feel free to open issues on this github repo.
