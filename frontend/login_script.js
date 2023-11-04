
// JavaScript to toggle between login and sign-up forms
const loginForm = document.querySelector('.container');
const signupForm = document.getElementById('signupForm');
const loginLink = document.getElementById('loginLink');
const signupLink = document.getElementById('signupLink');

loginLink.addEventListener('click', () => {
    loginForm.style.display = 'block';
    signupForm.style.display = 'none';
});

signupLink.addEventListener('click', () => {
    loginForm.style.display = 'none';
    signupForm.style.display = 'block';
});

