<?php
// Only inject the login payload if we are not already authenticated/redirected.
// Adminer adds the username to the GET parameters after successful login.
if (!isset($_GET['username'])) {
    $_POST['auth']['driver'] = 'pgsql';
    $_POST['auth']['server'] = getenv('ADMINER_SERVER') ?: 'db';
    $_POST['auth']['username'] = getenv('ADMINER_USER');
    $_POST['auth']['password'] = getenv('ADMINER_PASSWORD');
    $_POST['auth']['db'] = getenv('ADMINER_DB');
}

// Load the actual Adminer core
require('adminer.php');
?>
