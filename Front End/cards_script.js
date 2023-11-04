const swipeCard = document.getElementById('swipeCard');
let startX;
let currentCardIndex = 0;

const cardsData = [
    { name: 'Product 1', imageUrl: 'product1.jpg', description: 'Description 1' },
    { name: 'Product 2', imageUrl: 'product2.jpg', description: 'Description 2' },
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
        <p>${currentCardData.description}</p>
    `;
}

function handleTouchStart(event) {
    startX = event.touches[0].clientX;
}

function handleTouchMove(event) {
    if (!startX) return;

    const currentX = event.touches[0].clientX;
    const deltaX = currentX - startX;

    swipeCard.style.transform = `translateX(${deltaX}px)`;
}

function handleTouchEnd() {
    const threshold = 100; // Adjust the threshold as needed

    if (startX && startX - event.changedTouches[0].clientX > threshold) {
        // Swipe left, show the next card
        currentCardIndex = (currentCardIndex + 1) % cardsData.length;
    }

    swipeCard.style.transform = 'translateX(0)';
    startX = null;

    // Update the UI with the next card
    updateCard();
}

