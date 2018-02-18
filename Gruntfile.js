module.exports = function(grunt) {

    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        sass: {
            dist: {
                files: [{
                    expand: true,
                    cwd: 'sass',
                    src: ['*.scss'],
                    dest: 'public/css',
                    ext: '.css'
                }]
            }
        },
        cssmin: {
            target: {
                files: [{
                    expand: true,
                    cwd: 'css',
                    src: ['*.css', '!*.min.css'],
                    dest: 'public/css/min',
                    ext: '.min.css'
                }]
            }
        },
        watch: {
            sass: {
                files: ['sass/*.scss'],
                tasks: ['sass']
            },
            styles: {
                files: ['public/css/*.css'],
                tasks: ['cssmin']
            }
        }
    });
  
    grunt.loadNpmTasks('grunt-contrib-sass');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-watch');
  
    grunt.registerTask('default', ['sass', 'cssmin']);
  };