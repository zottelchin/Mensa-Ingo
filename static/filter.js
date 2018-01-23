function updateFilter(id, fv, fn) {
  fv = fv.split(/[^a-z0-9]+/i).map(x => x.trim().toLowerCase()).filter(x => x);
  if (!fv.length) {
    [...document.getElementsByClassName("meal")].forEach(e => e.classList.remove("hide-filter-" + id));
    return localStorage.removeItem("filter-" + id);
  }
  localStorage.setItem("filter-" + id, fv.join(","));
  [...document.getElementsByClassName("meal")].forEach(e => {
    if (!fn(e, fv)) {
      e.classList.add("hide-filter-" + id);
    } else {
      e.classList.remove("hide-filter-" + id);
    }
  })
}

function updateFilterS() {
  var fv = [...document.querySelectorAll(".symbol-filter [data-name]")].filter(e => !e.classList.contains("disabled")).map(e => e.getAttribute("data-name")).join(",");
  if (!document.querySelectorAll(".symbol-filter .disabled").length) fv = "";
  else fv += ",suppe";
  
  updateFilter("s", fv, (e, fv) =>
    [...e.querySelectorAll(".icons img")].filter(icon => {
      if (fv.indexOf(icon.src.match(/\/mensasym_(.*)\.png$/)[1]) > -1) return true;
      else return false;
    }).length
  );
}
function updateFilterA() {
  updateFilter("a", document.getElementById("filter-a").value, (e, fv) => {
    var x = e.querySelector(".hints").textContent.toLowerCase();
    return !fv.filter(f => x.indexOf("(" + f + ")") > -1).length;
  });
}

function symbolFilterAll() {
    [...document.querySelectorAll(".symbol-filter [data-name]")].forEach(e => e.classList.remove("disabled"));
    updateFilterS();
}
function symbolFilter(fv) {
    fv = fv.split(",");
    [...document.querySelectorAll(".symbol-filter [data-name]")].forEach(e => fv.indexOf(e.getAttribute("data-name")) > -1 ? e.classList.remove("disabled") : e.classList.add("disabled"));
    updateFilterS();
}
[...document.querySelectorAll(".symbol-filter [data-name]")].forEach(e => e.addEventListener("click", () => { e.classList.toggle("disabled"); updateFilterS(); }));


// Restore
if (window.location.hash == "#v") {
  localStorage.setItem("filter-s", "vegetarisch,vegan,suppe")
} else if (window.location.hash == "#vv") {
  localStorage.setItem("filter-s", "vegan,suppe")
}

var fs = localStorage.getItem("filter-s");
if (fs) {
    [...document.querySelectorAll(".symbol-filter [data-name]")].forEach(e => e.classList.add("disabled"));
    fs.split(",").filter(x => x && x != "suppe").forEach(x => document.querySelector(".symbol-filter [data-name=\"" + x + "\"]").classList.remove("disabled"));
}
document.getElementById("filter-a").value = localStorage.getItem("filter-a");
  
updateFilterS();
updateFilterA();
