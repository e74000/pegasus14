const termsContainer = document.querySelector('.terms-container');
const terms = document.querySelector('.terms');

termsContainer.addEventListener('scroll', () => {
    if (termsContainer.scrollTop + termsContainer.clientHeight >= terms.offsetHeight) {
        // When the user reaches the end of the terms, append the content again.
        terms.innerHTML += terms.innerHTML;
    }
});
