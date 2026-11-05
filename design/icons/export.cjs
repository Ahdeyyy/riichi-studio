const fs=require('fs'),path=require('path');
const {Resvg}=require(path.join(process.cwd(),'.cache/icon-tools/node_modules/@resvg/resvg-js'));
function png(svg,size){return new Resvg(svg,{fitTo:{mode:'width',value:size},font:{loadSystemFonts:false}}).render().asPng()}
function exportIcon(svg,sizeVariants={}){
 const out=path.join(__dirname,'final');fs.mkdirSync(out,{recursive:true});
 fs.writeFileSync('resources/icon.svg',svg);fs.writeFileSync(path.join(out,'icon.svg'),svg);
 const sizes=[16,20,24,32,40,48,64,128,256,512,1024];const images=new Map();
 for(const size of sizes){const image=png(sizeVariants[size]||svg,size);images.set(size,image);fs.writeFileSync(path.join(out,`icon-${size}.png`),image)}
 fs.writeFileSync('resources/icon.png',images.get(1024));
 const icoSizes=sizes.filter(n=>n<=256),header=Buffer.alloc(6+icoSizes.length*16);header.writeUInt16LE(1,2);header.writeUInt16LE(icoSizes.length,4);let offset=header.length;
 icoSizes.forEach((size,i)=>{const p=6+i*16;header[p]=size===256?0:size;header[p+1]=size===256?0:size;header.writeUInt16LE(1,p+4);header.writeUInt16LE(32,p+6);header.writeUInt32LE(images.get(size).length,p+8);header.writeUInt32LE(offset,p+12);offset+=images.get(size).length});
 fs.writeFileSync('resources/icon.ico',Buffer.concat([header,...icoSizes.map(s=>images.get(s))]));
 return images;
}
module.exports={png,exportIcon};
