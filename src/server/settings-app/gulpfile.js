var gulp = require('gulp')
var msx = require('gulp-msx')
var less = require('gulp-less');
var minifyCSS = require('gulp-minify-css');
var rename = require("gulp-rename");
var browserify = require('browserify');

var uglify = require('gulp-uglify');

var source = require('vinyl-source-stream');
var buffer = require('vinyl-buffer');

var nodemon = require('nodemon');
var watch = require('gulp-watch');
var runSequence = require('run-sequence');
 
gulp.task('build-less', function() {
  return gulp.src('./src/less/main.less')
    .pipe(less())
    .pipe(rename('main.css'))
    .pipe(gulp.dest('./dist/css'));
})

gulp.task("copy-assets", function(){
    return gulp.src('./assets/**', {base: './assets'})
        .pipe(gulp.dest('./dist/assets/'));
});

gulp.task('build-js', function(done) {    
    var bundler = browserify('./src/js/main.js',{transform: ['mithrilify']}).bundle()
    
    return bundler
        .pipe(source('main.js'))
        .pipe(rename('main.js'))
        .pipe(gulp.dest('./dist/js'));
});

gulp.task('minify-js', function(){
    return gulp.src('./dist/js/main.js')
        .pipe(buffer())
        .pipe(uglify())
        .pipe(rename('main.min.js'))
        .pipe(gulp.dest('./dist/js'));
})

gulp.task('minify-css', function(){
    return gulp.src('./dist/css/main.css')
        .pipe(minifyCSS())
        .pipe(rename('main.min.css'))
        .pipe(gulp.dest('./dist/css'));
})

gulp.task('minify', function() {
    runSequence(['minify-css', 'minify-js'])
})

gulp.task('build', function(){
    runSequence(['build-js', 'build-less'], ['minify'])
});

gulp.task('watch', ['build'], function() {
  gulp.watch('./src/js/**/*.js', function(){
      runSequence(['build-js'], ['minify-js'])
  });
  gulp.watch('./src/less/**/*.less', function() {
      runSequence(['build-less'], ['minify-css'])
  });
});