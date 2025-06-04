import Foundation

struct Poll: Identifiable, Decodable {
    let id: String
    let name: String?
    let invite_code: String?
    let role: String?
    let members: [String]?
    let created_by: String?
    let created_at: String?
    let updated_at: String?
}
