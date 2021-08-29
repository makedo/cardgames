const HOST = '/api'; //@TODO get from env

export default async function api(path, options = {method: 'GET'}) {
    const response = await fetch(`${HOST}${path}`, {
        method: options.method,
        mode: 'cors', // no-cors, *cors, same-origin
        cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
        headers: {
          'Content-Type': 'application/json'
        }
      });

    if (!response.ok) {
        const error = await response.text();
        throw new Error(error);
      }
    
    return await response.json();
}
