import Foundation

struct PollOptionResult: Identifiable, Decodable {
    let option_id: String
    let option_name: String
    let vote_count: Int
    let voter_ids: [String]

    var id: String { option_id }

    private enum CodingKeys: String, CodingKey {
        case option_id, option_name, vote_count, voter_ids
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.container(keyedBy: CodingKeys.self)
        option_id = try container.decode(String.self, forKey: .option_id)
        option_name = try container.decode(String.self, forKey: .option_name)
        vote_count = try container.decode(Int.self, forKey: .vote_count)
        voter_ids = try container.decodeIfPresent([String].self, forKey: .voter_ids) ?? []
    }
}
