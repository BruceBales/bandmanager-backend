# Band Management Platform

## Short-term Goals

    - Band system allowing users to form/join "bands"
        - Band Calendar allowing users to mark days as unavailable or available
        - Each band has a calendar, but calendars will be synchronized for users who are in multiple bands.
          For instance, if you are in two bands, each band will have it's own calendar, but if you have
          a show booked for one band, you will automatically be marked as "busy" on the other band's calendar.
          Essentially, there will be one calendar with multiple contextual views.
        - I will likely try to find a 3rd-party calendar platform to use for this, since
          writing my own calendar platform would be an extremely difficult task.
          My criteria will be that it's a calendar system that can be self-hosted, has a REST api,
          and is easy to embedd into a web page.
    - Venue system allowing users to create venues, and link other members as staff
    - Show system allowing bands and venues to book "shows" that show up on the calendar

## Long-term Goals

    - Project Management feature
        - File storage for demos, STEMS, chord sheets, ect.
        - Version management for project files
        - Embedded tab editor/player(this one is a long-shot that probably won't happen)
            - Maybe use TuxGuitar in the browser? I know it's in Java. Licensing might be a problem.
    -Collaboration marketplace


## Design Goals

    - Absolutely all logic relating to the service itself must happen in the Go Backend
        - I might split the backend into multiple services, but overall, it is critical that everything
          that happens is an API call to the backend API, so that both a website and mobile app can
          be maintained as separate projects.
        - The only programming that happens in the frontend components should be presentation of data. For
          example, showing which bands the active user is in should use an API call to a "getbands" endpoint
          that the user ID is passed to, rather than the UI doing a database lookup.
        - Theoretically, it should be fully possible to use the app as a user with nothing but a REST client making
          API calls.
    - Both the backend and frontend should be containerized
        - Backend binaries should be capable of running outside of container, so keep it simple. The purpose of
          the container is to make deployments safer- but the application should be capable of executing
          as one single binary per service.
    - Multiple instances of both backend and frontend services should be able to run at once
        - Memcached and Redis will be useful for this.
        - I will likely never need this, but I want to be able to just so I can feel cool.

## Plan of Attack

    - Build user system with sessions that expire.
        - Backend endpoint that validates login info, returns session ID
        - Session ID eventually runs out of time, forcing user to login again
        - Login endpoint takes credential input, returns session ID
        - Other endpoints use session ID to match user info
        - Session ID's shouldn't be visible to the user. PHP will have it's own session management that
          will correlate to the backend session ID, but will never make it available
          to a browser.
    - Build endpoints for band-related actions
    - Build endpoints for venue-related actions
