import { Kafka } from 'kafkajs';
import { readFileSync } from 'fs';
import avro, { Type } from 'avsc';
import commander from 'commander';

interface Message {
    [key: string]: any;
}

const kafka = new Kafka({
    clientId: 'my-app',
    brokers: ['localhost:9092']
});

const program = commander
    .option('-s, --schema <path>', 'Path to Avro schema file')
    .option('-p, --payload <path>', 'Path to payload JSON file')
    .option('-t, --topic <name>', 'Name of Kafka topic to produce to')
    .parse(process.argv);

const producer = kafka.producer();

if (!program.schema || !program.payload || !program.topic) {
    console.error('Missing required options');
    process.exit(1);
}

const schema = avro.Type.forSchema(readFileSync(program.schema));

const payload = JSON.parse(readFileSync(program.payload, 'utf-8')) as Message;

const encodedMessage = schema.toBuffer(payload);

async function sendToKafka(topic: string, message: Buffer) {
    await producer.connect();
    await producer.send({
        topic,
        messages: [
            { value: message }
        ]
    });
    await producer.disconnect();
}

sendToKafka(program.topic, encodedMessage)
    .then(() => console.log('Message sent to Kafka'))
    .catch(console.error);
