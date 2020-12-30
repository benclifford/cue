FROM debian:buster

RUN apt-get update && apt-get install -y sudo

RUN apt-get update && apt-get install -y certbot

# run this container like this:
# cue -D "-v /home/benc/etc-letsencrypt:/etc/letsencrypt" certbot sudo -i certbot certonly --manual -d '*.cqx.ltd.uk'
# (certbot renew doesnt work with manual/DNS mode for some reason, but certonly figures out it needs to renew...)
# (and there will be some manual DNS interaction needed which is why --manual is needed)

#IMPORTANT NOTES:
# - Congratulations! Your certificate and chain have been saved at:
#   /etc/letsencrypt/live/cqx.ltd.uk/fullchain.pem
#   Your key file has been saved at:
#   /etc/letsencrypt/live/cqx.ltd.uk/privkey.pem
#   Your cert will expire on 2019-11-13. To obtain a new or tweaked
#   version of this certificate in the future, simply run certbot
#   again. To non-interactively renew *all* of your certificates, run
#   "certbot renew"
# - Your account credentials have been saved in your Certbot
#   configuration directory at /etc/letsencrypt. You should make a
#   secure backup of this folder now. This configuration directory will
#   also contain certificates and private keys obtained by Certbot so
#   making regular backups of this folder is ideal.


# then on smtp-out container, place /etc/letsencrypt/live/cqx.ltd.uk/privkey.pem into /etc/mail/tls/sendmail-common.key
# and /etc/letsencrypt/live/cqx.ltd.uk/fullchain.pem into both sendmail-server.crt and sendmail-client.crt

# in srv-imap container, install new keys in /etc/courier/imapd.pem
#   first privkey.pem then fullchain.pem pasted together into that file

