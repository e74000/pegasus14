const video = document.getElementById("video-background");
video.playbackRate = 0.5;

function getSwipes() {
    fetch("/swipes/").then((response) => {
       return response.json()
    }).then((num) => {
        document.getElementById("swipes").textContent = String(num).padStart(9, "0").replace(/\B(?=(\d{3})+(?!\d))/g, ",") + " swipes"
    })
}

getSwipes()

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

function setLoginBasket() {
    login_token = getCookieAsJSON("login_token")

    const xhr = new XMLHttpRequest();
    xhr.open("PUT", "/validate_token/", true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4) {
            if (xhr.status >= 200 && xhr.status < 300) {
                console.log("Registration response:", xhr.responseText);

                document.getElementById("login-basket").innerHTML = ``
                document.getElementById("login-basket").innerText = ``



                document.getElementById("login-basket").innerHTML = `
                    <li><a href="/basket/">
                    <span class="material-symbols-outlined">
                        account_circle
                    </span>
                    </a>
                </li>`
            }
        }
    };

    xhr.send(JSON.stringify(login_token));
}

setLoginBasket()