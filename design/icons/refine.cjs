const fs=require('fs'),path=require('path');
const {svg,rect,circle,ivory,red,dir,Resvg}=require('./explore.cjs');
const {png}=require('./export.cjs');
const variants=[
 {id:'A',name:'Table declaration',angle:0,width:380,height:64,radius:9,dot:22,edge:7},
 {id:'B',name:'Clean diagonal',angle:-30,width:380,height:64,radius:9,dot:22,edge:0},
 {id:'C',name:'Placed stick',angle:-24,width:380,height:68,radius:8,dot:24,edge:6},
 {id:'D',name:'Quiet horizontal',angle:0,width:360,height:76,radius:6,dot:26,edge:0}
];
function mark(v,size=1024){
 const w=size<=32?384:v.width,h=size<=16?88:size<=32?76:v.height,d=size<=16?30:size<=32?26:v.dot;
 let m=v.edge?rect(256-w/2,256-h/2+v.edge,w,h,v.radius,'#D9C9A3'):'';
 m+=rect(256-w/2,256-h/2,w,h,v.radius,ivory);
 m+=size<=16?rect(256-d,256-d,d*2,d*2,6,red):circle(256,256,d,red);
 return `<g transform="rotate(${v.angle} 256 256)">${m}</g>`;
}
fs.mkdirSync(path.join(dir,'refinements'),{recursive:true});
for(const v of variants){for(const size of [16,32,64,256,1024]){const s=svg(mark(v,size));fs.writeFileSync(path.join(dir,'refinements',v.id+'-'+size+'.svg'),s);fs.writeFileSync(path.join(dir,'refinements',v.id+'-'+size+'.png'),png(s,size))}}
let s='<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="840"><rect width="100%" height="100%" fill="#E8E6DF"/>';
variants.forEach((v,i)=>{const y=20+i*205;s+=`<image x="20" y="${y}" width="176" height="176" href="data:image/png;base64,${fs.readFileSync(path.join(dir,'refinements',v.id+'-256.png')).toString('base64')}"/><text x="220" y="${y+28}" font-family="Segoe UI" font-size="20" fill="#233A32">${v.id}. ${v.name}</text>`;
 for(let mode=0;mode<2;mode++){const x0=mode?710:220;s+=rect(x0,y+52,460,120,12,mode?'#192A25':'#FAFAF7');let x=x0+20;for(const size of [16,32,64]){s+=`<image x="${x}" y="${y+77}" width="${size}" height="${size}" href="data:image/png;base64,${fs.readFileSync(path.join(dir,'refinements',v.id+'-'+size+'.png')).toString('base64')}"/><text x="${x}" y="${y+155}" font-family="Segoe UI" font-size="13" fill="${mode?'#E8E6DF':'#33443C'}">${size}px</text>`;x+=size+60}}
});s+='</svg>';fs.writeFileSync(path.join(dir,'round-2.svg'),s);fs.writeFileSync(path.join(dir,'round-2.png'),new Resvg(s,{font:{loadSystemFonts:false,fontFiles:['C:/Windows/Fonts/segoeui.ttf']}}).render().asPng());
fs.writeFileSync(path.join(dir,'refinements.json'),JSON.stringify(variants,null,2));
module.exports={mark,variants};
