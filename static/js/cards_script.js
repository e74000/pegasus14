document.addEventListener('DOMContentLoaded', async () => {
    const swipeCard = document.getElementById('swipeCard');
    let startX;
    let currentCard = null;

    function getCookieAsJSON(name) {
        const cookies = document.cookie.split(';');
        for (let i = 0; i < cookies.length; i++) {
            const cookie = cookies[i].trim();
            if (cookie.startsWith(name + '=')) {
                const cookieValue = cookie.substring(name.length + 1);
                try {
                    return JSON.parse(decodeURIComponent(cookieValue));
                } catch (error) {
                    console.error("Error parsing JSON from the cookie:", error);
                }
            }
        }
    }

    function out_of_products() {
        swipeCard.remove();
        console.log("Deleted the card");

        document.getElementById("out_of_products").innerText = "Out of Products!"
    }

    const login_token = getCookieAsJSON("login_token");
    const cards = [];

    async function initCard() {
        try {
            await fillCards();
            updateCard();
        } catch (error) {
            console.error("Error initializing card: ", error);
        }

        swipeCard.addEventListener('touchstart', handleTouchStart, false);
        swipeCard.addEventListener('touchmove', handleTouchMove, false);
        swipeCard.addEventListener('touchend', handleTouchEnd, false);
    }

    function updateCard() {
        if (cards.length > 0) {
            currentCard = cards.pop()[0];
            console.log(currentCard);

            swipeCard.innerHTML = `
        <h1 class="cardHeader">${currentCard.title}</h1>
        <img class="cardImage" src="${currentCard.img}" alt="Product Image">
        <div class="cardDescription">
            <span>${currentCard.description}</span>
            <span>Price: <b>${currentCard.price}</b></span>
        </div>
      `;
        } else {
            out_of_products();
        }
    }

    async function fillCards() {
        try {
            const response = await fetch(`/suggest/${login_token.email}`);
            if (!response.ok) {
                throw new Error("Network response was not ok"); // Fixed "throw an Error" to "throw new Error"
            }
    
            const data = await response.json();
    
            for (const sku of data) {
                const skuApiUrl = `/product/${sku}`;
                const productResponse = await fetch(skuApiUrl);
    
                if (!productResponse.ok) {
                    throw new Error("Network response was not ok");
                }
    
                const product = await productResponse.json();
                cards.push(product);
            }
        } catch (error) {
            console.error("Error in fillCards: ", error);
        }
    }
    

    function handleTouchStart(event) {
        startX = event.touches[0].clientX;
    }

    function handleTouchMove(event) {
        const threshold = 100;
    
        if (!startX) return;
    
        const currentX = event.touches[0].clientX;
        const deltaX = 1.2 * (currentX - startX);
    
        const intensity = 0.2 * Math.pow(Math.min(1.2, Math.abs(deltaX) / threshold), 2);
    
        let red, green, blue;
    
        if (deltaX > 0) {
            // Gradually turn green as it moves to the right
            red = 255 - Math.round(255 * intensity);
            green = 255;
            blue = 255 - Math.round(255 * intensity);
        } else {
            // Gradually turn red as it moves to the left
            red = 255;
            green = 255 - Math.round(255 * intensity);
            blue = 255 - Math.round(255 * intensity);
        }
    
        swipeCard.style.backgroundColor = `rgba(${red}, ${green}, ${blue}, 0.8)`;
        swipeCard.style.transform = `translateX(${deltaX}px)`;
    }
    
    function handleTouchEnd(event) {
        const threshold = 100;

        if (startX) {
            const deltaX = startX - event.changedTouches[0].clientX;

            if (deltaX > threshold || deltaX < -threshold) {
                let impression = {
                    email: login_token.email,
                    sku: currentCard.sku,
                    swipe: 0 + (deltaX > 0),
                    claim: login_token,
                };

                updateCard();

                fetch("/impression/", {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(impression),
                })
                    .then((data) => {
                        // Handle the response data if needed
                        console.log("PUT request successful:", data);
                    })
                    .catch((error) => {
                        console.error("Error sending PUT request: ", error);
                    });

                if (deltaX < -threshold) {
                    showAddedToBasketPopup();
                }
            }
        }

        // Reset styles
        swipeCard.style.transform = 'translateX(0)';
        swipeCard.style.backgroundColor = ''; // Reset background color

        startX = null;
    }

    function closeAddedToBasketModal() {
        const modal = document.getElementById('addedToBasketModal');
        modal.style.display = 'none';
    }
    
    function showAddedToBasketPopup() {
        const modal = document.getElementById('addedToBasketModal');
        modal.style.display = 'flex';
    
        // Close the modal after 0.5 seconds
        setTimeout(() => {
            closeAddedToBasketModal();
        }, 500); // Changed the timeout to 500ms for half a second
    }

    initCard();
});
