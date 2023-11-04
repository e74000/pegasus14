const swipeCard = document.getElementById('swipeCard');
let startX;
let currentCardIndex = 0;

const cardsData = [
    { name: 'Product 1', imageUrl: 'product1.jpg', description: 'Description 1', price: '$19.99' },
    { name: 'Product 2', imageUrl: 'product2.jpg', description: 'Description 2', price: '$29.99' },
    { name: 'Product 3', imageUrl: 'product3.jpg', description: 'Description 3', price: '$39.99' },
    { name: 'Product 4', imageUrl: 'product4.jpg', description: 'Description 4', price: '$49.99' },
    { name: 'Product 5', imageUrl: 'product5.jpg', description: 'Description 5', price: '$59.99' },
    // Add more card data as needed
];

initCard();

function initCard() {
    updateCard();
    swipeCard.addEventListener('touchstart', handleTouchStart, false);
    swipeCard.addEventListener('touchmove', handleTouchMove, false);
    swipeCard.addEventListener('touchend', handleTouchEnd, false);
}

function updateCard() {
    const currentCardData = cardsData[currentCardIndex];
    swipeCard.innerHTML = `
    <h1 class="cardHeader">${currentCardData.name}</h1>
    <img class="cardImage" src="${currentCardData.imageUrl}" alt="Product Image">
    <div class="cardDescription">
        <span>${currentCardData.description}</span>
        <span>Price: <b>${currentCardData.price}</b></span>
    </div>
    `;
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
        // Gradually turn gray as it moves to the left
        const grayIntensity = 0.5 * Math.pow(Math.min(1.2, Math.abs(deltaX) / threshold), 2);
        red = green = blue = 255 - Math.round(127 * grayIntensity);
    }

    swipeCard.style.backgroundColor = `rgba(${red}, ${green}, ${blue},0.8)`;
    swipeCard.style.transform = `translateX(${deltaX}px)`;
}




function handleTouchEnd(event) {
    const threshold = 100; // Adjust the threshold as needed

    if (startX) {
        const deltaX = startX - event.changedTouches[0].clientX;

        if (deltaX > threshold) {
            // Swipe left, show the next card
            currentCardIndex = (currentCardIndex + 1) % cardsData.length;
        } else if (deltaX < -threshold) {
            // Swipe right, show the previous card
            currentCardIndex = (currentCardIndex - 1 + cardsData.length) % cardsData.length;

            // Show a popup for "Added to Basket"
            showAddedToBasketPopup();
        }
    }

    // Reset styles
    swipeCard.style.transform = 'translateX(0)';
    swipeCard.style.backgroundColor = ''; // Reset background color

    startX = null;

    // Update the UI with the next card
    updateCard();
}

// Your existing JavaScript code

function showAddedToBasketPopup() {
    const modal = document.getElementById('addedToBasketModal');
    modal.style.display = 'flex';

    // Close the modal after 0.5 seconds
    setTimeout(() => {
        closeAddedToBasketModal();
    }, 300);
}

function closeAddedToBasketModal() {
    const modal = document.getElementById('addedToBasketModal');
    modal.style.display = 'none';
}