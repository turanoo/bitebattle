import Foundation

struct Restaurant: Identifiable, Decodable {
    let place_id: String
    let name: String
    let address: String?
    let rating: Double?
    let photos: [Photo]?

    struct Photo: Decodable {
        let photo_reference: String
    }

    var id: String { place_id }
}
