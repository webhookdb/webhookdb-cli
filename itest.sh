#!/usr/bin/env bash

export WEBHOOKDB_API_HOST=http://localhost:18001
export WEBHOOKDB_LOG_LEVEL=debug
export EXE=./webhookdb

make build

echo "*** Hello and welcome to the webhookdb integration test."

echo "*** testing auth logout..."
${EXE} auth logout


echo "*** Please enter a primary email to use for auth:"
read primaryemail
${EXE} auth login --username=$primaryemail
echo "*** Great! Now check that primary email or DB (customer_reset_codes) and enter the OTP here: "
read primaryotp
${EXE} auth otp --username=$primaryemail --token=$primaryotp
echo "*** testing services list..."
${EXE} services list

echo "*** testing org create..."
${EXE} org create
echo "What is the key for the the new org? i.e. what is the slugified version of the org name you just typed in:"
read orgkey

echo "*** testing org list..."
${EXE} org list

echo "*** testing org remove and optional org flag..."
${EXE} org invite --org=$orgkey --username=albertobalsam@aphex.co
${EXE} org remove --org=$orgkey --username=albertobalsam@aphex.co

echo "*** testing org activate..."
${EXE} org activate $orgkey

echo "*** testing org current..."
${EXE} org current


echo "*** Please enter a secondary email to test org invite functionality:"
read secondaryemail
echo "*** testing org invite..."
${EXE} org invite --username=$secondaryemail

echo "*** testing org changerole..."
${EXE} org invite --username=albertobalsam@aphex.co
${EXE} org changerole --usernames=$secondaryemail,albertobalsam@aphex.co --role=admin

${EXE} auth logout
${EXE} auth login --username=$secondaryemail

echo "*** Great! We're gonna log into the second account in order to test join functionality and finish the test"
echo "*** Now check that secondary email and enter the OTP here: "
read secondaryotp
${EXE} auth otp --username=$secondaryemail --token=$secondaryotp


echo "*** testing org join..."
echo "*** In your secondary email inbox, you should have also received an email with a join code. Entwher that join code here:"
read joincode
${EXE} org join $joincode

echo "*** testing org members..."
${EXE} org members

echo "*** testing integrations create..."
${EXE} integrations create fake_v1
echo "What is the opaque_id of the new integration?"
read opaqueId

echo "*** testing integrations list..."
${EXE} integrations list
echo "*** What is the table name listed above?"
read tableName

echo "*** testing backfill..."
${EXE} backfill $opaqueId

echo "*** testing db tables..."
${EXE} db tables

echo "*** testing db sql..."
${EXE} db sql "SELECT * FROM $tableName"

echo "*** testing subscription info..."
${EXE} subscription info

echo "*** testing subscription edit..."
${EXE} subscription edit