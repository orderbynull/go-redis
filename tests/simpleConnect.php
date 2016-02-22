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

//for ($i = 0; $i < 100; $i++) {
//    $client->ping();
//}

//$client->ping();

$key = "mysimplekey";
$value = "example";

$client->set($key, $value);

$returnValue = $client->get($key);
var_dump(["RETURN", $returnValue]);
assert($returnValue == $value);
