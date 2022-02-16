source = ["./dist/webhookdb-macos-amd64_darwin_amd64/webhookdb"]
bundle_id = "com.lithictech.webhookdbcli"

apple_id {
  username = "matt@lithic.tech"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Lithic Technology LLC (V8WYJK5CN8)"
}