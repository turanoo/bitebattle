import SwiftUI
import CoreLocation

struct PollOptionView: View {
    let poll: Poll
    @Environment(\.presentationMode) var presentationMode

    @State private var zipCode: String = ""
    @State private var selectedCategory: String = ""
    @State private var isSearching: Bool = false
    @State private var statusMessage: String?
    @State private var searchResults: [Restaurant] = []
    @State private var selectedRestaurants: Set<String> = []
    @State private var isSubmitting: Bool = false

    // Example categories
    let categories = ["Pizza", "Sushi", "Italian", "Ramen", "Burgers", "Mexican", "Chinese", "Indian", "Thai"]

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

    var body: some View {
        ZStack {
            LinearGradient(
                gradient: Gradient(colors: [Color.pink.opacity(0.7), Color.orange.opacity(0.7)]),
                startPoint: .topLeading,
                endPoint: .bottomTrailing
            )
            .ignoresSafeArea()

            VStack(spacing: 20) {
                Text("Add Option to \(poll.name)")
                    .font(.title2)
                    .fontWeight(.bold)
                    .foregroundColor(.white)
                    .padding(.top, 16)

                VStack(spacing: 12) {
                    TextField("Enter Zip Code", text: $zipCode)
                        .keyboardType(.numbersAndPunctuation)
                        .padding()
                        .background(Color.white.opacity(0.9))
                        .cornerRadius(10)
                        .shadow(radius: 1)

                    Picker("Category", selection: $selectedCategory) {
                        Text("Select Category").tag("")
                        ForEach(categories, id: \.self) { cat in
                            Text(cat).tag(cat)
                        }
                    }
                    .pickerStyle(MenuPickerStyle())
                    .padding()
                    .background(Color.white.opacity(0.9))
                    .cornerRadius(10)
                    .shadow(radius: 1)

                    Button(action: {
                        searchRestaurants()
                    }) {
                        HStack {
                            Image(systemName: "magnifyingglass")
                                .foregroundColor(.white)
                            Text("Search")
                                .fontWeight(.semibold)
                                .foregroundColor(.white)
                        }
                        .frame(maxWidth: .infinity, minHeight: 44)
                        .padding()
                        .background((zipCode.isEmpty || selectedCategory.isEmpty) ? Color.gray.opacity(0.5) : Color.pink.opacity(0.8))
                        .cornerRadius(12)
                        .shadow(radius: 2)
                    }
                    .disabled(zipCode.isEmpty || selectedCategory.isEmpty || isSearching)
                }

                if isSearching {
                    ProgressView()
                        .progressViewStyle(CircularProgressViewStyle(tint: .white))
                        .padding()
                } else if let status = statusMessage {
                    Text(status)
                        .foregroundColor(.white)
                        .font(.footnote)
                        .padding(.top, 8)
                        .shadow(radius: 1)
                }

                ScrollView {
                    VStack(spacing: 16) {
                        ForEach(searchResults) { restaurant in
                            Button(action: {
                                toggleSelection(for: restaurant.place_id)
                            }) {
                                HStack(alignment: .top, spacing: 12) {
                                    if let photoRef = restaurant.photos?.first?.photo_reference {
                                        AsyncImage(url: photoURL(photoRef)) { image in
                                            image.resizable()
                                        } placeholder: {
                                            Color.gray.opacity(0.3)
                                        }
                                        .frame(width: 60, height: 60)
                                        .cornerRadius(8)
                                    } else {
                                        Color.gray.opacity(0.3)
                                            .frame(width: 60, height: 60)
                                            .cornerRadius(8)
                                    }
                                    VStack(alignment: .leading, spacing: 4) {
                                        Text(restaurant.name)
                                            .font(.headline)
                                            .foregroundColor(.pink)
                                        if let rating = restaurant.rating {
                                            Text("Rating: \(String(format: "%.1f", rating))")
                                                .font(.subheadline)
                                                .foregroundColor(.white)
                                        }
                                        if let address = restaurant.address, !address.isEmpty {
                                            Text(address)
                                                .font(.caption)
                                                .foregroundColor(.white)
                                        }
                                    }
                                    Spacer()
                                    if selectedRestaurants.contains(restaurant.place_id) {
                                        Image(systemName: "checkmark.circle.fill")
                                            .foregroundColor(.green)
                                            .font(.title2)
                                    }
                                }
                                .padding()
                                .background(selectedRestaurants.contains(restaurant.place_id) ? Color.green.opacity(0.15) : Color.white.opacity(0.15))
                                .cornerRadius(14)
                                .shadow(radius: 2)
                            }
                        }
                    }
                    .padding(.horizontal, 8)
}

                if !searchResults.isEmpty {
                    Button(action: {
                        submitSelectedOptions()
                    }) {
                        HStack {
                            Image(systemName: "paperplane.fill")
                                .foregroundColor(.white)
                            Text(isSubmitting ? "Submitting..." : "Submit")
                                .fontWeight(.semibold)
                                .foregroundColor(.white)
                        }
                        .frame(maxWidth: .infinity, minHeight: 44)
                        .padding()
                        .background(selectedRestaurants.isEmpty || isSubmitting ? Color.gray.opacity(0.5) : Color.orange.opacity(0.8))
                        .cornerRadius(12)
                        .shadow(radius: 2)
                    }
                    .disabled(selectedRestaurants.isEmpty || isSubmitting)
                    .padding(.top, 8)
                }

                Spacer()
            }
            .padding()
        }
        .navigationTitle("Add Option")
        .navigationBarTitleDisplayMode(.inline)
    }

    func searchRestaurants() {
        guard !zipCode.isEmpty, !selectedCategory.isEmpty else { return }
        statusMessage = nil
        isSearching = true
        searchResults = []
        selectedRestaurants = []

        // Geocode zip code to coordinates
        let geocoder = CLGeocoder()
        geocoder.geocodeAddressString(zipCode) { placemarks, error in
            if let error = error {
                DispatchQueue.main.async {
                    self.statusMessage = "Invalid zip code: \(error.localizedDescription)"
                    self.isSearching = false
                }
                return
            }
            guard let location = placemarks?.first?.location else {
                DispatchQueue.main.async {
                    self.statusMessage = "Could not find location for zip code."
                    self.isSearching = false
                }
                return
            }
            let coordinates = "\(location.coordinate.latitude),\(location.coordinate.longitude)"
            fetchRestaurants(category: selectedCategory, location: coordinates)
        }
    }

    func fetchRestaurants(category: String, location: String) {
        guard let token = UserDefaults.standard.string(forKey: "authToken"),
              !token.isEmpty,
              let url = URL(string: "http://localhost:8080/api/restaurants/search?q=\(category.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? "")&location=\(location.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed) ?? "")") else {
            statusMessage = "Not logged in."
            isSearching = false
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")

        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                isSearching = false
                if let error = error {
                    statusMessage = "Failed: \(error.localizedDescription)"
                    return
                }
                guard let data = data else {
                    statusMessage = "No data returned."
                    return
                }
                do {
                    let decoded = try JSONDecoder().decode([Restaurant].self, from: data)
                    self.searchResults = decoded
                } catch {
                    statusMessage = "Failed to load restaurants."
                }
            }
        }.resume()
    }

    func toggleSelection(for placeId: String) {
        if selectedRestaurants.contains(placeId) {
            selectedRestaurants.remove(placeId)
        } else {
            selectedRestaurants.insert(placeId)
        }
    }

    func submitSelectedOptions() {
        guard !selectedRestaurants.isEmpty else { return }
        isSubmitting = true
        statusMessage = nil

        let selected = searchResults.filter { selectedRestaurants.contains($0.place_id) }

        // Prepare the array payload
        let payload: [[String: Any]] = selected.map { restaurant in
            let imageUrl = restaurant.photos?.first.flatMap { photoURL($0.photo_reference)?.absoluteString } ?? ""
            return [
                "restaurant_id": restaurant.place_id,
                "name": restaurant.name,
                "image_url": imageUrl as String,
                "menu_url": "" // You can add menu URL if available
            ]
        }

        guard let token = UserDefaults.standard.string(forKey: "authToken"),
            !token.isEmpty,
            let url = URL(string: "http://localhost:8080/api/polls/\(poll.id)/options"),
            let httpBody = try? JSONSerialization.data(withJSONObject: payload) else {
            statusMessage = "Not logged in or failed to encode data."
            isSubmitting = false
            return
        }

        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = httpBody

        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                isSubmitting = false
                if let error = error {
                    statusMessage = "Failed: \(error.localizedDescription)"
                    return
                }
                presentationMode.wrappedValue.dismiss()
            }
        }.resume()
    }


    func photoURL(_ reference: String) -> URL? {
        // Replace with your backend's photo endpoint if needed
        // This is a placeholder for Google Places API photo endpoint
        let maxWidth = 400
        return URL(string: "https://maps.googleapis.com/maps/api/place/photo?maxwidth=\(maxWidth)&photoreference=\(reference)&key=AIzaSyCHFSr3LuANYJu4dALcGZquzF0MX8OSKsI")
    }
}
