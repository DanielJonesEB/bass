const styleKey = "theme";
const linkId = "theme";

const controlsId = "choosetheme"
const switcherId = "styleswitcher";
const resetId = "resetstyle";

function storeStyle(style) {
  window.localStorage.setItem(styleKey, style);
}

function loadStyle() {
  return window.localStorage.getItem(styleKey);
}

function setActiveStyle(style) {
  var link = document.getElementById(linkId);
  if (link) {
    link.href = "css/base16/base16-"+style+".css";
  } else {
    link = document.createElement('link');
    link.id = linkId;
    link.rel = "stylesheet";
    link.type = "text/css";
    link.href = "css/base16/base16-"+style+".css";
    link.media = "all";
    document.head.appendChild(link);
  }

  var switcher = document.getElementById(switcherId);
  if (switcher) {
    // might not be loaded yet; this function is called twice, once super early
    // to prevent flickering, and again once all the dom is loaded up
    switcher.value = style;
  }

  resetReset();
}

function resetReset() {
  var style = loadStyle();
  var reset = document.getElementById(resetId);
  if (!style) {
    if (reset) {
      // no style selected; remove reset element
      reset.remove();
    }

    return
  }

  if (reset) {
    // no style and no reset; done
    return
  }

  // has style but no reset element
  reset = document.createElement("a");
  reset.id = resetId;
  reset.onclick = resetStyle;
  reset.href = 'javascript:void(0)';
  reset.text = "reset";
  reset.className = "reset";

  var chooser = document.getElementById(controlsId);
  if (chooser) {
    chooser.prepend(reset);
  }
}

function setStyleOrDefault(def) {
  setActiveStyle(loadStyle() || def);
}

function switchStyle(event) {
  var style = event.target.value;
  storeStyle(style);
  setActiveStyle(style);
}

function resetStyle() {
  window.localStorage.removeItem(styleKey);
  setActiveStyle(defaultStyle);
}

// dark-mode media query matched or not
var defaultDarkStyle = "rose-pine";
var defaultLightStyle = "rose-pine-dawn";
var defaultStyle = defaultLightStyle;

if (window.matchMedia) {
  let prefersDark = '(prefers-color-scheme: dark)'
  let dark = window.matchMedia(prefersDark).matches;

  defaultStyle = dark ? defaultDarkStyle : defaultLightStyle;

  window.matchMedia(prefersDark).addEventListener('change', event => {
    if (event.matches) {
      defaultStyle = defaultDarkStyle;
    } else {
      defaultStyle = defaultLightStyle;
    }

    setStyleOrDefault(defaultStyle);
  })
}

setStyleOrDefault(defaultStyle);

window.onload = function() {
  // call again to update switcher selection
  setStyleOrDefault(defaultStyle);
}
