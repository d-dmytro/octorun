const { randomBytes } = require("crypto");
const { faker } = require("@faker-js/faker");

const INTERVAL = 1000;

const command = process.argv[2];

const getRandomInt = (max) => Math.floor(Math.random() * max);

if (command === "crash" && index === 10) {
  // Force the process to exit with error
  undefined.hello;
}

if (command === "waitandexit") {
  setTimeout(() => {
    process.stdout.write(`Exiting...\n`);
  }, 5000);
  return;
}

let index = 0;

setInterval(() => {
  index += 1;
  const outputStream = getRandomInt(2) === 0 ? "stdout" : "stderr";
  const message = faker.lorem.sentence();
  process[outputStream].write(`${message}\n`);
}, INTERVAL);
