const fs = require('fs');
const path = require('path');

// Get the environment variables
const googleMapsApiKey = process.env.GOOGLE_MAPS_API_KEY || '';

// Path to the generated environment files
const targetPathProd = path.resolve(__dirname, 'src/environments/environment.prod.js');
const targetPathDev = path.resolve(__dirname, 'src/environments/environment.js');

// Read the file contents and replace the placeholders
function replaceInFile(filePath) {
  if (!fs.existsSync(filePath)) {
    console.log(`File does not exist: ${filePath}`);
    return;
  }

  let fileContent = fs.readFileSync(filePath, 'utf8');
  
  // Replace the Google Maps API key placeholder
  fileContent = fileContent.replace(/googleMapsApiKey: ['"]GOOGLE_MAPS_API_KEY['"]/, `googleMapsApiKey: '${googleMapsApiKey}'`);
  
  // Write the modified content back to the file
  fs.writeFileSync(filePath, fileContent, 'utf8');
  
  console.log(`Environment values replaced in ${filePath}`);
}

// Check if the production environment file exists and replace values
if (fs.existsSync(targetPathProd)) {
  replaceInFile(targetPathProd);
}

// Check if the development environment file exists and replace values
if (fs.existsSync(targetPathDev)) {
  replaceInFile(targetPathDev);
}

console.log('Environment setup completed');