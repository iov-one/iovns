DOMAIN="hello"
ORIGINAL_OWNER=$(iovnscli keys show fd -a)
STATIC_OWNER=$(iovnscli keys show dp -a)
NEW_OWNER=$(iovnscli keys show ok -a)

echo registering domain owned by ${ORIGINAL_OWNER}
yes | iovnscli tx starname register-domain --domain ${DOMAIN} --from fd --type closed
sleep 5
echo register account from ${ORIGINAL_OWNER} owned by ${STATIC_OWNER}
yes | iovnscli tx starname register-account --name "never-change" --domain hello --from fd --owner ${STATIC_OWNER}
sleep 5
echo register account that should change from ${ORIGINAL_OWNER} owned by ${ORIGINAL_OWNER}
yes | iovnscli tx starname register-account --name "should-change" --domain hello --from fd
sleep 5
yes | iovnscli tx starname set-account-metadata --domain hello --name fd --metadata "Why the uri suffix?" --from fd
sleep 5
yes | iovnscli tx starname set-account-metadata --domain hello --name "" --metadata "something" --from fd
sleep 5
echo transferring domain to ${NEW_OWNER}
yes | iovnscli tx starname transfer-domain --domain hello --from fd --new-owner ${NEW_OWNER} --transfer-flag 1
sleep 5
echo the following account should be owned by  ${STATIC_OWNER}
iovnscli query starname resolve --domain hello --name never-change
echo the following account should be owned by ${NEW_OWNER}
iovnscli query starname resolve --domain hello --name should-change
