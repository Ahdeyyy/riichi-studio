const fs=require('fs'),path=require('path');
const {rect,dir,Resvg}=require('./explore.cjs');
const {png,exportIcon}=require('./export.cjs');
const master=fs.readFileSync(path.join(dir,'refinements/A-1024.svg'),'utf8');
function snapped(size){
 const x=Math.round(size/8),w=size-2*x,h=Math.round(size/8),y=Math.floor(size/2-h/2),d=size<=24?2:4,r=size<=24?0:1;
 return `<svg xmlns="http://www.w3.org/2000/svg" width="${size}" height="${size}" viewBox="0 0 ${size} ${size}">${rect(size*24/512,size*24/512,size*464/512,size*464/512,size*105/512,'#173F36')}${rect(x,y+h,w,1,r,'#D9C9A3')}${rect(x,y,w,h,r,'#F5EDDA')}${rect(Math.floor((size-d)/2),y+Math.floor((h-d)/2),d,d,size<=24?0:1,'#CB493B')}</svg>`;
}
const sizes={};for(const n of [16,20,24,32,40,48])sizes[n]=snapped(n);
const images=exportIcon(master,sizes);
fs.writeFileSync('resources/icon-window.png',png(snapped(32),128));
for(const [n,s] of Object.entries(sizes))fs.writeFileSync(path.join(dir,'final',`icon-${n}.svg`),s);
let proof='<svg xmlns="http://www.w3.org/2000/svg" width="1000" height="500"><rect width="100%" height="100%" fill="#E8E6DF"/>';
proof+=`<image x="24" y="24" width="360" height="360" href="data:image/png;base64,${images.get(512).toString('base64')}"/>`;
proof+='<text x="24" y="440" font-family="Segoe UI" font-size="20" fill="#233A32">Table declaration</text>';
for(let mode=0;mode<2;mode++){const y=mode?270:24;proof+=rect(420,y,556,218,16,mode?'#192A25':'#FAFAF7');let x=440;for(const n of [16,24,32,48,64,128]){proof+=`<image x="${x}" y="${y+36}" width="${n}" height="${n}" href="data:image/png;base64,${images.get(n).toString('base64')}"/><text x="${x}" y="${y+190}" font-family="Segoe UI" font-size="13" fill="${mode?'#E8E6DF':'#33443C'}">${n}</text>`;x+=n+26}}
proof+='</svg>';fs.writeFileSync(path.join(dir,'final-proof.png'),new Resvg(proof,{font:{loadSystemFonts:false,fontFiles:['C:/Windows/Fonts/segoeui.ttf']}}).render().asPng());
