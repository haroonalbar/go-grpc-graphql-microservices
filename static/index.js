// Function to send GraphQL requests
function sendGraphQLRequest(query, targetElementId) {
  fetch('http://localhost:8000/graphql', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ query }),
  })
    .then(response => response.json())
    .then(data => {
      if (targetElementId === 'orders-response') {
        if (data.errors) {
          document.getElementById(targetElementId).innerHTML = `<p class="text-red-500">Error: ${data.errors[0].message}</p>`;
        } else {
          const accounts = data.data.accounts;
          let cardHTML = '';
          accounts.forEach(account => {
            cardHTML += `<div class="bg-white shadow-md rounded-lg p-6 mb-4">
                          <h3 class="text-xl font-semibold mb-2">${account.name}</h3>`;
            var sum = 0
            account.orders.forEach(order => {
              sum += order.totalPrice
            })
            // cardHTML += `<div class="bg-white shadow-md rounded-lg p-6 mb-4">
            //               <h3 class="text-xl font-semibold mb-2">${account.name}</h3>`;
            account.orders.forEach(order => {
              sum += order.totalPrice
              cardHTML += `<div class="border-t border-gray-300 pt-4">
                            <!-- <h4 class="text-lg font-semibold">Total Price: $${order.totalPrice}</h4> -->
                            <ul class="list-disc pl-5">`;
              order.products.forEach(product => {
                cardHTML += `<li>${product.name} - ${product.description} (Quantity: ${product.quantity}, Price: $${product.price})</li>`;
              });
              cardHTML += `</ul></div>`;
            });
            var sum = 0
            account.orders.forEach(order => {
              sum += order.totalPrice
            })
            cardHTML += `<h3 class="text-xl font-semibold mb-2">Total Price : ${sum}</h3>`;
            cardHTML += `</div>`;
          });
          document.getElementById(targetElementId).innerHTML = cardHTML;
        }
      } else if (targetElementId === 'accounts-response') {
        const accounts = data.data.accounts;
        let tableHTML = '<table class="min-w-full border-collapse border border-gray-300"><thead><tr><th class="border border-gray-300 p-2">ID</th><th class="border border-gray-300 p-2">Name</th></tr></thead><tbody>';
        accounts.forEach(account => {
          tableHTML += `<tr><td class="border border-gray-300 p-2">${account.id}</td><td class="border border-gray-300 p-2">${account.name}</td></tr>`;
        });
        tableHTML += '</tbody></table>';
        document.getElementById(targetElementId).innerHTML = tableHTML;
      } else if (targetElementId === 'products-response') {
        if (data.errors) {
          document.getElementById(targetElementId).innerHTML = `<p class="text-red-500">Error: ${data.errors[0].message}</p>`;
        } else {
          const products = data.data.products;
          let tableHTML = '<table class="min-w-full border-collapse border border-gray-300"><thead><tr><th class="border border-gray-300 p-2">ID</th><th class="border border-gray-300 p-2">Name</th><th class="border border-gray-300 p-2">Description</th><th class="border border-gray-300 p-2">Price</th></tr></thead><tbody>';
          products.forEach(product => {
            tableHTML += `<tr><td class="border border-gray-300 p-2">${product.id}</td><td class="border border-gray-300 p-2">${product.name}</td><td class="border border-gray-300 p-2">${product.description}</td><td class="border border-gray-300 p-2">${product.price}</td></tr>`;
          });
          tableHTML += '</tbody></table>';
          document.getElementById(targetElementId).innerHTML = tableHTML;
        }
      } else {
        document.getElementById(targetElementId).innerHTML = JSON.stringify(data, null, 2);
      }
    })
    .catch(error => {
      console.error('Error:', error);
      document.getElementById(targetElementId).innerHTML = `<p class="text-red-500">An error occurred. Please try again.</p>`;
    });
}

// Create Account
document.getElementById('create-account-form').addEventListener('submit', function (event) {
  event.preventDefault();
  const accountName = document.getElementById('account-name').value;
  const query = `
                mutation {
                    createAccount(account: {name: "${accountName}"}) {
                        id
                        name
                    }
                }
            `;
  sendGraphQLRequest(query, 'account-response');
});

// Create Product
document.getElementById('create-product-form').addEventListener('submit', function (event) {
  event.preventDefault();
  const productName = document.getElementById('product-name').value;
  const productDescription = document.getElementById('product-description').value;
  const productPrice = document.getElementById('product-price').value;
  const query = `
                mutation {
                    createProduct(product: {name: "${productName}", description: "${productDescription}", price: ${productPrice}}) {
                        id
                        name
                        price
                    }
                }
            `;
  sendGraphQLRequest(query, 'product-response');
});

// Create Order
document.getElementById('create-order-form').addEventListener('submit', function (event) {
  event.preventDefault();
  const accountId = document.getElementById('order-account-id').value;
  const productId = document.getElementById('order-product-id').value;
  const quantity = document.getElementById('order-quantity').value;
  const query = `
                mutation {
                    createOrder(order: {accountId: "${accountId}", products: [{id: "${productId}", quantity: ${quantity}}]}) {
                        id
                        totalPrice
                        products {
                            name
                            quantity
                        }
                    }
                }
            `;
  sendGraphQLRequest(query, 'order-response');
});

// Query Accounts
document.getElementById('query-accounts').addEventListener('click', function () {
  const query = `
                query {
                    accounts {
                        id
                        name
                    }
                }
            `;
  sendGraphQLRequest(query, 'accounts-response');
});

// Query Products
document.getElementById('query-products').addEventListener('click', function () {
  const searchTerm = document.getElementById('product-search-term').value;
  const skip = parseInt(document.getElementById('product-skip').value) || 0;
  const take = parseInt(document.getElementById('product-take').value) || 5;

  const query = `
                query {
                    products(pagination: {skip: ${skip}, take: ${take}}, query: "${searchTerm}") {
                        id
                        name
                        description
                        price
                    }
                }
            `;
  sendGraphQLRequest(query, 'products-response');
});

// Query Orders 
document.getElementById('query-orders').addEventListener('click', function () {
  const accountId = document.getElementById('orders-account-id').value;
  const query = `
                query {
                  accounts(id: "${accountId}") {
                    name
                    orders {
                      id
                      totalPrice
                      products{
                        name
                        description
                        quantity
                        price
                      }
                    }
                  }
                }
            `;
  sendGraphQLRequest(query, 'orders-response');
});
