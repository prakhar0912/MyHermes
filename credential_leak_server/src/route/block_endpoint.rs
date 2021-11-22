use crate::utils::md5_encode;
use actix_web::{HttpResponse, Responder};
use redis::{AsyncCommands, RedisResult};
use std::ops::Add;

pub async fn block_endpoint(
    redis_client: actix_web::web::Data<redis::Client>,
    req: actix_web::HttpRequest,
) -> impl Responder {
    let endpoint = req.headers().get("endpoint");

    let endpoint = match endpoint {
        Some(v) => v.to_str().unwrap(),
        None => "",
    };
    let endpointv2 = md5_encode(endpoint.as_bytes());
    let conn = redis_client.get_async_connection().await;
    if let Ok(mut conn) = conn {
        let _: RedisResult<()> = conn
            .set(
                (String::from(endpointv2)).add("_blocked_endpoint"),
                String::from(endpoint),
            )
            .await;
        HttpResponse::Ok()
    } else {
        HttpResponse::InternalServerError()
    }
}
