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

login_token = getCookieAsJSON("login_token")

function loadProducts() {
    console.log("getting products")
    fetch(`/basket/${login_token.email}`).then((response) => {
        return response.json()
    }).then((data) => {
        let products = document.getElementById("products")
        data.forEach((sku) => {
            fetch(`/product/${sku}`)
                .then((product) => {
                    if (!product.ok) {
                        throw new Error("Network response was not ok")
                    }
                    return product.json()
                })
                .then((product) => {
                    let l = document.createElement("li")
                    l.innerText = product[0].title
                    products.appendChild(l)
                })
        })
    })
}

loadProducts()