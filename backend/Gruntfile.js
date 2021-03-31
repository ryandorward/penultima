var path = require('path');

module.exports = function (grunt) {
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        concurrent: {
            dev: {
                tasks: ['nodemon'],
                options: {
                    logConcurrentOutput: true
                }
            }
        },
        nodemon: {
            dev: {
                script: 'src/app.js',
                options: {
                    ignore: [
                        'node_modules/**'
                    ],
                    ext: 'js'
                }
            }
        },
        clean: {}
    });
    grunt.loadNpmTasks('grunt-concurrent');
    grunt.loadNpmTasks('grunt-nodemon');
    grunt.registerTask('default', ['concurrent']);
};
