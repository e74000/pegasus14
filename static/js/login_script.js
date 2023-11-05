
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

document.getElementById("registration-form").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent the default form submission

    // Get the values of email and password inputs
    const email = document.getElementById("register_email").value;
    const password = document.getElementById("register_password").value;

    // Create a JSON object
    const data = {
        "email": email,
        "password": password
    };

    // Convert the JSON object to a string
    const jsonData = JSON.stringify(data);

    // Send the JSON data to "/register/" using an AJAX request
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", "/register/", true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if (xhr.status >= 200 && xhr.status < 300) {
                // Handle the response here
                console.log("Registration response:", xhr.responseText);
                // You can add more code to handle the response as needed
            } else {
                console.error("Registration error:", xhr.status);
            }
        }
    };

    xhr.send(jsonData);
});

document.getElementById("validation-form").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent the default form submission

    // Get the values of email and password inputs
    const email = document.getElementById("validate_email").value;
    const password = document.getElementById("validate_password").value;

    // Create a JSON object
    const data = {
        "email": email,
        "password": password
    };

    // Convert the JSON object to a string
    const jsonData = JSON.stringify(data);

    // Send the JSON data to "/register/" using an AJAX request
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", "/validate/", true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if (xhr.status >= 200 && xhr.status < 300) {
                console.log("Registration response:", xhr.responseText);

                // Extract the login token from the response, assuming it's in a variable called token
                const token = xhr.responseText;

                // Save the token as a cookie with an expiry date (e.g., 1 day from now)
                const expirationDate = new Date();
                expirationDate.setDate(expirationDate.getDate() + 1);

                // Set the cookie
                document.cookie = `login_token=${token}; expires=${expirationDate.toUTCString()}; path=/`;

                location.href = "/app/"
            } else {
                console.error("Registration error:", xhr.status);
                document.getElementById("register_response_message").textContent = "incorrect username or password"
            }
        }
    };

    xhr.send(jsonData);
});