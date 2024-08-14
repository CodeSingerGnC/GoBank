import http from 'k6/http';
import { sleep } from 'k6';

// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
export const options = {
  vus: 10,
  duration: '30s',
};

export default function() {
  // Define the URL to which the POST request will be sent
  const url = 'http://localhost:8080/v1/user/login';

  // Define the payload to be sent in the POST request
  const payload = JSON.stringify({
    "user_account":"remotecodesingergnc",
    "password":"secret"
  });

  // Define the headers for the POST request
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  // Send the POST request
  const response = http.post(url, payload, params);

  // Log the response status code
  console.log(`Response status: ${response.status}`);

  // Sleep for 1 second
  sleep(1);
}