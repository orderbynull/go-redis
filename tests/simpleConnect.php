<?php
/**
 * @author Patsura Dmitry https://github.com/ovr <talk@dmtry.me>
 */

include_once __DIR__ . '/../vendor/autoload.php';

$client = new Predis\Client([
    'scheme' => 'tcp',
    'host'   => 'localhost',
    'port'   => 3333,
]);

for ($i = 0; $i < 100; $i++) {
    $client->ping();
}

