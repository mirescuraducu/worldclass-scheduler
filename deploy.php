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

task('go_build', function () {
    cd('{{deploy_path}}/current');
    run("go build initial.go");
    run("go mod vendor");
    run("go build classIds.go");
});


after('deploy:failed', 'deploy:unlock');
after('deploy:success', 'go_build');
