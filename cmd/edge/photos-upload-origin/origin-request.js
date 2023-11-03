exports.processOriginRequest = async (event, cb) => {


    const request = event.Records[0].cf.request;


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
    const extension = request.uri.split('.').pop().toLowerCase();

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

    request.uri = `/originals/${imageMeta.id}.${extension}`

    console.log("request :: ", JSON.stringify(request))


    return cb(null, request)
}