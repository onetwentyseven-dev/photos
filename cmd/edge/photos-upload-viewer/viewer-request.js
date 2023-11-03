const corsHeaders = {
    'access-control-allow-origin': [{
        key: 'Access-Control-Allow-Origin',
        value: 'https://photos.onetwentyseven.dev'
    }],
    "access-control-allow-credentials": [{
        "key": "Access-Control-Allow-Credentials",
        "value": "true"
    }],
    "access-control-allow-methods": [{
        "key": "Access-Control-Allow-Methods",
        "value": "OPTIONS,PUT"
    }],

}


exports.processViewerRequest = async (event, cb) => {

    const request = event.Records[0].cf.request;
    if (request.method === "OPTIONS") {
        return cb(null, {
            status: '204',
            statusDescription: "No Content",
        })
    }

    if (request.method !== "PUT") {
        return cb(null, {
            status: '405',
            statusDescription: 'Method not allowed',
            headers: {
                ...corsHeaders
            }
        })
    }

    return cb(null, request)

}