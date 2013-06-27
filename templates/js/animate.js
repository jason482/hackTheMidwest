$(document).ready(function () {

   var stripedKitty = $('#stripedKitty');
   var pig = $('#pig');
   var whiteKitty = $('#whiteKitty');
   var furry = $('#furry');
   var lizard = $('#lizard');
   var puppy = $('#puppy');
   var bunny = $('#bunny');
   var girlPuppy = $('#girlPuppy');

var resources = [
];
var symbols = {
"stage": {
   version: "2.0.0",
   minimumCompatibleVersion: "2.0.0",
   build: "2.0.0.250",
   baseState: "Base State",
   initialState: "Base State",
   gpuAccelerate: false,
   resizeInstances: false,
   content: {
         dom: [
],
         symbolInstances: [

         ]
      },
   states: {
      "Base State": {
         "#stripedKitty": [
            ["transform", "translateX", '-135.32%'],
            ["transform", "translateY", '33.76%']
         ],
         "#bunny": [
            ["transform", "translateX", '-108.22%'],
            ["transform", "translateY", '51.08%']
         ],
         "${_puppy}": [
            ["transform", "translateX", '1328.47%'],
            ["transform", "translateY", '2.54%']
         ],
         "${_girlPuppy}": [
            ["transform", "translateX", '-263.4%'],
            ["transform", "translateY", '249.12%']
         ],
         "${_lizard}": [
            ["transform", "translateX", '1014.93%'],
            ["transform", "translateY", '-79.8%']
         ],
         "${_furry}": [
            ["transform", "translateX", '2050.87%'],
            ["transform", "translateY", '472.08%']
         ],
         "${_whiteKitty}": [
            ["transform", "translateX", '925.97%'],
            ["transform", "translateY", '82.01%']
         ],
         "${_pig}": [
            ["transform", "translateX", '-162.67%'],
            ["transform", "translateY", '287.57%']
         ]
      }
   },
   timelines: {
      "Default Timeline": {
         fromState: "Base State",
         toState: "",
         duration: 1750,
         autoPlay: true,
         timeline: [
            { id: "eid33", tween: [ "transform", "${_furry}", "translateY", '379.91%', { fromValue: '472.08%'}], position: 1066, duration: 459, easing: "easeOutQuad" },
            { id: "eid28", tween: [ "transform", "${_whiteKitty}", "translateY", '163.07%', { fromValue: '82.01%'}], position: 83, duration: 982, easing: "easeInQuart" },
            { id: "eid14", tween: [ "transform", "${_puppy}", "translateY", '166.08%', { fromValue: '2.54%'}], position: 328, duration: 1422, easing: "easeInQuad" },
            { id: "eid12", tween: [ "transform", "${_puppy}", "translateX", '597.09%', { fromValue: '1328.47%'}], position: 328, duration: 1422, easing: "easeInQuad" },
            { id: "eid26", tween: [ "transform", "${_whiteKitty}", "translateX", '545.37%', { fromValue: '925.97%'}], position: 83, duration: 982, easing: "easeInQuart" },
            { id: "eid24", tween: [ "transform", "#stripedKitty", "translateY", '104.59%', { fromValue: '33.76%'}], position: 547, duration: 1040, easing: "easeOutCubic" },
            { id: "eid2", tween: [ "transform", "${_lizard}", "translateX", '685.79%', { fromValue: '1014.93%'}], position: 219, duration: 1424, easing: "easeOutQuart" },
            { id: "eid38", tween: [ "transform", "#bunny", "translateY", '103.65%', { fromValue: '51.08%'}], position: 0, duration: 852, easing: "easeInSine" },
            { id: "eid8", tween: [ "transform", "${_girlPuppy}", "translateX", '436.48%', { fromValue: '-263.4%'}], position: 455, duration: 1295, easing: "easeInCubic" },
            { id: "eid22", tween: [ "transform", "#stripedKitty", "translateX", '206.03%', { fromValue: '-135.32%'}], position: 547, duration: 1040, easing: "easeOutCubic" },
            { id: "eid20", tween: [ "transform", "${_pig}", "translateY", '193.8%', { fromValue: '287.57%'}], position: 0, duration: 1587, easing: "easeInQuint" },
            { id: "eid10", tween: [ "transform", "${_girlPuppy}", "translateY", '188.2%', { fromValue: '249.12%'}], position: 455, duration: 1295, easing: "easeInCubic" },
            { id: "eid18", tween: [ "transform", "${_pig}", "translateX", '162.72%', { fromValue: '-162.67%'}], position: 0, duration: 1587, easing: "easeInQuint" },
            { id: "eid37", tween: [ "transform", "#bunny", "translateX", '478.87%', { fromValue: '-108.22%'}], position: 0, duration: 852, easing: "easeInSine" },
            { id: "eid16", tween: [ "transform", "${_furry}", "translateX", '1390.01%', { fromValue: '2050.87%'}], position: 1066, duration: 459, easing: "easeOutQuad" },
            { id: "eid4", tween: [ "transform", "${_lizard}", "translateY", '186.11%', { fromValue: '-79.8%'}], position: 219, duration: 1424, easing: "easeOutQuart" }         ]
      }
   }
}
}

});