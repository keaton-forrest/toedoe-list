// script.js


function init() {

    const modal = document.getElementById('connection-error');

    let polling = false;

    // Handle connection error events where the client can't connect to the server
    document.addEventListener('htmx:sendError', function(x) {
    
        if (polling) return;
        
        polling = true;

        console.log(x); // What does this event look like?

        // Adding this class will make the modal visible and prevent the user from interacting with the page
        modal.classList.add('htmx-request');
    
        // Keep checking the server for a connection
        const id = setInterval(() => {
            fetch('/ping')
                .then(response => {
                    if (response.status === 200) {
                        polling = false;
                        clearInterval(id);
                        modal.classList.remove('htmx-request');
                    }
                })
                .catch(err => console.error(err));
        }, 1000);
    });
}

init();