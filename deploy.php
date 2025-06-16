<?php
namespace Deployer;

require 'recipe/common.php';

// Config

set('repository', 'git@github.com:mirescuraducu/worldclass-scheduler.git');

add('shared_files', []);
add('shared_dirs', []);
add('writable_dirs', []);

// Hosts

host('alina-si-radu.ro')
    ->set('remote_user', 'deployer')
    ->set('deploy_path', '~/worldclass-scheduler');

// Hooks

after('deploy:failed', 'deploy:unlock');
