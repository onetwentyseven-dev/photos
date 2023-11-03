import {
    S3Client,
    GetObjectCommand,
    PutObjectCommand,
} from "@aws-sdk/client-s3";


const l = (m) => console.log(m);

const sharp = require('sharp');

const apiHost = "https://photos.onetwentyseven.dev"

exports.processOriginResponse = async (event, cb) => {
    l(`processOriginResponse :: ${JSON.stringify(event)}`)


    const s3Client = new S3Client();

    const response = event.Records[0].cf.response;

    // if (response.status != 404) {
    //     return cb(null, response)
    // }

    const request = event.Records[0].cf.request;
    // Check to see if the x-photos-bucket header is set on the reques 
    if (!request.headers['x-photos-bucket']) {
        l("No x-photos-bucket header set on request")
        return cb(null, {
            status: '400',
            statusDescription: 'Bad Request',
            headers: {
                ...response.headers,
            }
        })
    }

    const bucket = request.headers['x-photos-bucket'][0].value;

    // trim '/originals/' and the extensions from the uri
    const name = request.uri.replace(/^\/originals\//, '').split('.')[0];

    await s3Client.send(
        new GetObjectCommand({
            Bucket: bucket,
            Key: request.uri,
        })
    ).then(data => {
        return sharp(data.Body).resize(200, 200).toFormat('jpeg').toBuffer()
    }).then(buffer => {
        s3Client.send(
            new PutObjectCommand({
                Bucket: bucket,
                Key: `/thumbnails/${name}.jpeg`,
                Body: buffer,
                ContentType: 'image/jpeg',
            })
        ).then(() => {

            l("Image Resized Successfully and Uploaded")
            l("Updating Image Status to Processed")
            return patchImageStatus(request, "processed")

        }).catch(e => {

            l("Error uploading resized image :: ", e)
            l(e)
            return patchImageStatus(request, "errored", e)

        })
    }).catch(e => {
        l("Error resizing image :: ", e)
        l(e)
        return patchImageStatus(request, "errored", e)
    });


    // await s3.getObject({
    //     Bucket: bucket,
    //     Key: request.uri,
    // }).promise().then(data => {
    //     return sharp(data.Body).resize(200, 200).toFormat('jpeg').toBuffer()
    // }).then(buffer => {
    //     s3.putObject({
    //         Bucket: bucket,
    //         Key: `/thumbnails/${name}.jpeg`,
    //         Body: buffer,
    //         ContentType: 'image/jpeg',
    //     }).promise().then(() => {

    //         l("Image Resized Successfully and Uploaded")
    //         l("Updating Image Status to Processed")
    //         return patchImageStatus(request, "processed")

    //     }).catch(e => {

    //         l("Error uploading resized image :: ", e)
    //         return patchImageStatus(request, "errored", e)

    //     })
    // }).catch(e => {
    //     l("Error resizing image :: ", e)
    //     return patchImageStatus(request, "errored", e)
    // });

    return cb(null, response);
}


async function patchImageStatus(request, status, e) {

    const name = request.uri.replace(/^\/originals\//, '').split('.')[0];

    const body = {
        "id": name,
        "status": status
    }
    if (e) {
        body.error = e.toString()
    }

    const patchOptions = {
        method: "PATCH",
        headers: {
            'cookie': request.headers.cookie[0].value
        },
        body: JSON.stringify(body)
    }

    return fetch(`${apiHost}/api/image/metadata`, patchOptions).then(r => {
        if (r.status !== 200) { throw new Error("Invalid auth cookie") }
        l("Image Status Updated Successfully")
    }).catch(e => {
        l("Error updating image status :: ", e)
        l(e)
        return e;
    })

}