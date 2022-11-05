

# Design



What are we trying to build?




Load Ballancer ->
    web server ->
        /.build (starts a build if one isn't running; displays build logs)
        /* (returns static content)



Developer mode:

  A docker image that generates the static system to a known directory.

Production mode:

  A docker image that responds to /.build and returns static data otherwise.




Alternative design:

Somehow get 
  admin.polyfriendly.org -> goes to port 80 (the panel service)
  www.polyfriendly.org -> reads static data

How would that work?  A shared volume that is updated by the admin page but shared
statically by the hosting company?



/.build
  Template based on /usr/local/www/www.scrollodex.net/html/kickit/main-hugo-bi.html 
/.build/logtail.js
  Template based on kickit/logtail-main-hugo-bi.js

/logtail.css
Static file


