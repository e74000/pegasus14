// JavaScript code to infinitely repeat and increase corruption intensity when reaching the bottom
const contentContainer = document.getElementsByClassName("terms-container")[0];
const contentToRepeat = document.getElementsByClassName("terms")[0];
let corruptionIntensity = 1;

function repeatAndCorruptContent() {
    const clone = contentToRepeat.cloneNode(true);
    const paragraphs = clone.getElementsByTagName("p");

    for (const paragraph of paragraphs) {
        // Apply text corruption with intensity
        paragraph.textContent = corruptText(paragraph.textContent, corruptionIntensity);
    }

    contentContainer.appendChild(clone);

    // Increase corruption intensity for the next iteration
    corruptionIntensity += 1;
}

function corruptText(text, intensity) {
    // Apply a more intense text corruption effect based on the intensity
    // You can customize this to create a more complex corruption effect
    return text.split('').map(char => {
        if (Math.random() < (intensity * 0.1)) {
            // Replace characters with random characters
            return String.fromCharCode(32 + Math.floor(Math.random() * 2 * corruptionIntensity * 1000));
        }
        return char;
    }).join('');
}

function isAtBottom() {
    return window.innerHeight + window.scrollY >= document.body.offsetHeight;
}

window.addEventListener("scroll", () => {
    if (isAtBottom()) {
        repeatAndCorruptContent();
    }
});
