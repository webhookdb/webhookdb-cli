#!/usr/bin/env bash

alias whtest="API_HOST=http://localhost:18001 ./webhookdb"

make build


echo "*** Hello and welcome to the webhookdb integration test."

echo "\n*** testing auth logout..."
whtest auth logout


echo "\n*** Please enter a primary email to use for auth:"
read primaryemail
whtest auth login --username=$primaryemail
echo "*** Great! Now check that primary email and enter the OTP here: "
read primaryotp
whtest auth otp --username=$primaryemail --token=$primaryotp
echo "\n*** testing services list..."
whtest services list

echo "\n*** testing org create..."
whtest org create
echo "What is the key for the the new org? i.e. what is the slugified version of the org name you just typed in:"
read orgkey

echo "\n*** testing org list..."
whtest org list

echo "\n*** testing org remove and optional org flag..."
whtest org invite --org=$orgkey --username=albertobalsam@aphex.co
whtest org remove --org=$orgkey --username=albertobalsam@aphex.co

echo "\n*** testing org activate..."
whtest org activate $orgkey

echo "\n*** testing org current..."
whtest org current


echo "\n*** Please enter a secondary email to test org invite functionality:"
read secondaryemail
echo "*** testing org invite..."
whtest org invite --username=$secondaryemail

echo "\n*** testing org changerole..."
whtest org invite --username=albertobalsam@aphex.co
whtest org changerole --usernames=$secondaryemail,albertobalsam@aphex.co --role=admin

whtest auth logout
whtest auth login --username=$secondaryemail

echo "\n*** Great! We're gonna log into the second account in order to test join functionality and finish the test"
echo "*** Now check that secondary email and enter the OTP here: "
read secondaryotp
whtest auth otp --username=$secondaryemail --token=$secondaryotp


echo "\n*** testing org join..."
echo "*** In your secondary email inbox, you should have also received an email with a join code. Entwher that join code here:"
read joincode
whtest org join $joincode

echo "\n*** testing org members..."
whtest org members

echo "\n*** testing integrations create..."
whtest integrations create fake_v1
echo "What is the opaque_id of the new integration?"
read opaqueId

echo "\n*** testing integrations list..."
whtest integrations list
echo "*** What is the table name listed above?"
read tableName

echo "\n*** testing backfill..."
whtest backfill $opaqueId

echo "\n*** testing db tables..."
whtest db tables

echo "\n*** testing db sql..."
whtest db sql "SELECT * FROM $tableName"

echo "\n*** testing subscription info..."
whtest subscription info

echo "\n*** testing subscription edit..."
whtest subscription edit