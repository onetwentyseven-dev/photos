

const hooks = {
    'viewer-request': handleViewerRequest,
    // 'origin-request': handleOriginRequest
}

exports.handler = async (
    event, _, callback
) => {
    const config = event.Records[0].cf.config;
    const eventType = config.eventType

    console.log("EVENT", eventType, JSON.stringify(event))

    const hook = hooks[eventType];
    if (hook) {
        return hook(event, callback);
    }

    switch (eventType) {
        case 'viewer-request':
        case 'origin-request':
            console.log('case viewer-request | origin-request', event.Records[0].cf.request.headers)
            callback(null, event.Records[0].cf.request);
            return
        case 'viewer-response':
        case 'origin-response':
            console.log('case viewer-response | origin-response', event.Records[0].cf.response)
            callback(null, event.Records[0].cf.response);
            return
    }

}

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

async function handleViewerRequest(event, cb) {

    console.log("VIEWER REQUEST EVENT", JSON.stringify(event))

    const request = event.Records[0].cf.request;
    if (request.method === "OPTIONS") {

        const response = {
            status: '204',
            statusDescription: "No Content",
        }

        return cb(null, response)
    }

    if (request.method === "PUT") {
        console.log(JSON.stringify(request))

        let isValid = undefined;

        const fetchOptions = {
            method: "GET",
            headers: {
                'cookie': request.headers.cookie[0].value
            }
        }

        console.log("fetchOptions :: ", JSON.stringify(fetchOptions))

        await fetch("https://photos.onetwentyseven.dev/api/auth/validate", fetchOptions)
            .then(r => {
                console.log("Response status code :: ", r.status);
                if (r.status !== 200) { throw new Error("Invalid auth cookie") }

                return r.json();
            })
            .then(r => {
                isValid = r;
            })
            .catch(e => {
                console.log(`Failed to validate auth cookie`, e.message, e.response);
            });

        if (!isValid) {
            return cb(null, {
                status: '401',
                statusDescription: 'Unauthorized',
                headers: {
                    ...corsHeaders
                }
            })
        }

        // Trim leading slash from the beginning of the URI
        const name = request.uri.replace(/^\//, '');
        const extension = request.uri.split('.').pop();

        const postOptions = {
            method: "POST",
            headers: {
                'cookie': request.headers.cookie[0].value
            },
            body: JSON.stringify({
                "name": name.toLowerCase(),
            })

        }

        console.log("POST OPTIONS :: ", JSON.stringify(postOptions))

        let imageMeta = undefined

        await fetch("https://photos.onetwentyseven.dev/api/image/metadata", postOptions).then(r => {
            console.log("Response status code :: ", r.status);
            if (r.status !== 200) { throw new Error("Invalid auth cookie") }
            return r.json();
        }).then(r => {
            imageMeta = r;
            return r;
        }).catch(e => {
            console.log(`Failed to validate auth cookie`, e.message, e.response);
        });

        console.log("IMAGE META :: ", imageMeta)


        request.headers['x-amz-meta-user-id'] = [{
            key: 'x-amz-meta-user-id',
            value: imageMeta.user_id
        }]

        // Parse the reuqest.uri for the file extension.
        // ex: /images/image.jpg


        request.uri = `/${imageMeta.id}.${extension}`

        console.log("request :: ", JSON.stringify(request))


        return cb(null, request)
    }

    return cb(null, {
        status: '405',
        statusDescription: 'Method not allowed',
        headers: {
            ...corsHeaders
        }
    })
}

// async function handleOriginRequest(event, cb) {

//     console.log("ORIGIN REQUEST EVENT", JSON.stringify(event))

//     const request = event.Records[0].cf.request;

//     const origin = request.origin.s3

//     origin.customHeaders = {
//         ...origin.customHeaders,
//         'x-amz-meta-image-id': [{
//             key: 'x-amz-meta-image-id',
//             value: 'test'
//         }]
//     }

//     console.log("origin :: ", JSON.stringify(origin))
//     console.log("request :: ", JSON.stringify(request))

//     return cb(null, request)

// }