import http from 'k6/http';
import { check } from 'k6';

export const options = {
  stages: [
    { duration: "10s", target: 5 },
    { duration: "45s", target: 100 },
    { duration: "15s", target: 0 },
  ],
};

// The function that defines VU logic.
//
// See https://grafana.com/docs/k6/latest/examples/get-started-with-k6/ to learn more
// about authoring k6 scripts.
//
export default function() {
  const payload = JSON.stringify({
    source: 'delivery',
    dishes: [
      { name: (Math.random() + 1).toString(36).substring(7) },
    ],
    time: new Date().toISOString(),
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post('http://127.0.0.1:9001/api/v1/orders', payload, params);
  check(res, { 'status was 201': (r) => r.status == 201 });
}
